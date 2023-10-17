using MR.Net.Proto;
using MR.Net.Proto.Battle;
using System;
using System.Collections.Generic;
using System.Reflection.Metadata;
using System.Threading.Tasks;

namespace MR.BattleServer {
    public class Room {
        private const int COUNT_NUM = 5;

        private int m_ID;
        private int m_UdpPortMin;
        private int m_updPortMax;
        private List<User> m_Users = new List<User>();
        private Game m_Game;
        private bool m_Running = true;
        private bool m_ReadyRelease;
        private bool m_Counting;
        private DateTime m_BeginTime;

        public RoomPB PB {
            get {
                var pb = new RoomPB();
                lock (m_Users) {
                    pb.Players = new PlayerPB[m_Users.Count];
                    for (int i = 0; i < m_Users.Count; i++)
                        pb.Players[i] = m_Users[i].PB;
                }
                return pb;
            }
        }

        public Room(int id, int udpPortMin, int udpPortMax) {
            m_ID = id;
            m_UdpPortMin = udpPortMin;
            m_updPortMax = udpPortMax;
            RoomLoop();
        }

        private async void RoomLoop() {
            while (m_Running) {
                if (m_Game != null)
                    await GameLoop();
                else
                    await ReadyLoop();
            }
        }

        private async Task ReadyLoop() {
            while (!m_Counting || DateTime.Now < m_BeginTime)
                await Task.Delay(100);
            var players = new TcpProtoHandle[6];
            var watchers = new List<TcpProtoHandle>();
            for (int i = 0; i < m_Users.Count; i++) {
                var user = m_Users[i];
                if (user.Seat < 0)
                    watchers.Add(user.Handle);
                else
                    players[user.Seat] = user.Handle;
            }
            m_Game = new Game(players, watchers, m_UdpPortMin, m_updPortMax);
        }


        private async Task GameLoop() {
            while (m_Game != null) {
                m_Game.Update();
                if (m_Game.GameState == Game.State.Begin) {
                    SendAll(new GameRunS2C());
                    break;
                }
                await Task.Delay(10);
            }
            while (m_Game != null && m_Game.GameState == Game.State.Begin) {
                m_Game.Update();
                await Task.Delay(10);
            }
            if (m_Game != null) {
                foreach (var u in m_Users)
                    u.Ready = false;
                m_Counting = false;
                SendAll(new GameEndS2C { Sync = m_Game.Sync });
                SendAll(new RoomDataS2C { RoomData = PB });
                m_Game.Dispose();
                m_Game = null;
            }
            if (m_ReadyRelease)
                ReleaseRoom();
        }

        public void Enter(TcpProtoHandle handle, string name, ulong[] weapons) {
            var user = new User(this, handle, name, weapons);
            SendAll(new PlayerUpdateS2C { Mode = UpdatePBType.Add, Player = user.PB });
            m_Users.Add(user);
            handle.Send(new EnterRoomS2C { RoomID = m_ID });
            handle.Send(new RoomDataS2C { RoomData = PB });
        }

        private void SendAll<T>(T obj) {
            foreach (var user in m_Users)
                user.Handle.Send(obj);
        }

        private CodePBType UserExit(User user) {
            if (!m_Users.Contains(user))
                return CodePBType.Success;
            if (m_Users.IndexOf(user) > 0) {
                m_Users.Remove(user);
                SendAll(new PlayerUpdateS2C { Mode = UpdatePBType.Remove, Player = user.PB });
                if (m_Counting && user.Seat != -1) {
                    m_Counting = false;
                    SendAll(new GameCountCancelS2C());
                }
            } else if (m_Game != null) {
                m_Users.Remove(user);
                SendAll(new PlayerUpdateS2C { Mode = UpdatePBType.Remove, Player = user.PB });
                m_ReadyRelease = true;
            } else {
                m_Users.Remove(user);
                ReleaseRoom();
            }
            return CodePBType.Success;
        }

        private void ReleaseRoom() {
            m_Running = false;
            SendAll(new RoomBreakupS2C());
            BattleServer.Instance.RoomRelease(m_ID);
        }

        private CodePBType SitDown(User user, int seat) {
            if (m_Game != null)
                return CodePBType.GameInProgress;
            if (m_Counting)
                return CodePBType.GameInCount;
            lock (m_Users) {
                for (int i = 0; i < m_Users.Count; i++) {
                    var p = m_Users[i];
                    if (p.Seat == seat)
                        return CodePBType.SeatOccupied;
                }
            }
            user.Seat = seat;
            SendAll(new PlayerUpdateS2C { Mode = UpdatePBType.Update, Player = user.PB });
            return CodePBType.Success;
        }

        public CodePBType StandUp(User user) {
            if (m_Game != null)
                return CodePBType.GameInProgress;
            user.Seat = -1;
            user.Ready = false;
            SendAll(new PlayerUpdateS2C { Mode = UpdatePBType.Update, Player = user.PB });
            if (m_Counting) {
                m_Counting = false;
                SendAll(new GameCountCancelS2C());
            }
            return CodePBType.Success;
        }

        public CodePBType Move(User user, int seat) {
            if (m_Game != null)
                return CodePBType.GameInProgress;
            if (user.Ready)
                return CodePBType.ReadyCanNotMove;
            lock (m_Users) {
                for (int i = 0; i < m_Users.Count; i++) {
                    var p = m_Users[i];
                    if (p.Seat == seat)
                        return CodePBType.SeatOccupied;
                }
            }
            user.Seat = seat;
            SendAll(new PlayerUpdateS2C { Mode = UpdatePBType.Update, Player = user.PB });
            return CodePBType.Success;
        }

        public CodePBType Ready(User user, bool value) {
            if (user.Ready == value)
                return CodePBType.Success;
            if (user.Seat == -1)
                return CodePBType.ReadyNeadSit;
            user.Ready = value;
            SendAll(new PlayerUpdateS2C { Mode = UpdatePBType.Update, Player = user.PB });
            if (m_Counting && !value) {
                m_Counting = false;
                SendAll(new GameCountCancelS2C());
            }
            return CodePBType.Success;
        }

        public CodePBType Count(User user, bool start) {
            if (start && !m_Counting) {
                if (user != m_Users[0])
                    return CodePBType.NoPermission;
                if (!m_Users.Exists(m => m.Seat != -1))
                    return CodePBType.PlayerNumberError;
                foreach (var u in m_Users)
                    if (u != m_Users[0] && u.Seat != -1 && !u.Ready)
                        return CodePBType.NotAllReady;
                m_Counting = true;
                m_BeginTime = DateTime.Now + TimeSpan.FromSeconds(COUNT_NUM);
                SendAll(new GameCountStartS2C { Scend = COUNT_NUM });
            }
            if (!start && m_Counting) {
                m_Counting = false;
                SendAll(new GameCountCancelS2C());
            }
            return CodePBType.Success;
        }

        public CodePBType GameResult(User user, ScorePB[] score) {
            if (m_Game == null)
                return CodePBType.GameNotBegin;
            m_Game.ReciveScore(user.Handle, score);
            return CodePBType.Success;
        }

        public class User {
            private Room m_Room;
            public TcpProtoHandle Handle { get; private set; }
            public int Seat { get; set; } = -1;
            public bool Ready { get; set; }
            private string m_Name;
            private ulong m_Character;
            private ulong[] m_Weapons;

            public PlayerPB PB {
                get {
                    var pb = new PlayerPB();
                    pb.Name = m_Name;
                    pb.Seat = Seat;
                    if (Seat != -1) {
                        pb.Character = m_Character;
                        pb.Weapons = m_Weapons;
                        pb.Ready = Ready;
                    }
                    return pb;
                }
            }

            public User(Room room, TcpProtoHandle handle, string name, ulong[] weapons) {
                m_Room = room;
                Handle = handle;
                m_Name = name;
                m_Weapons = weapons;
                Regist();
            }
            private void Regist() {
                Handle.onClose += UnRegist;
                Handle.Regist<ExitRoomC2S>(OnExitRoom);
                Handle.Regist<SitDownC2S>(OnSitDown);
                Handle.Regist<StandUpC2S>(OnStandUp);
                Handle.Regist<MoveC2S>(OnMove);
                Handle.Regist<ReadyC2S>(OnReady);
                Handle.Regist<CountC2S>(OnCount);
                Handle.Regist<GameResultC2S>(OnGameResult);
            }

            private void UnRegist() {
                Handle.onClose -= UnRegist;
                Handle.UnRegist<ExitRoomC2S>(OnExitRoom);
                Handle.UnRegist<SitDownC2S>(OnSitDown);
                Handle.UnRegist<StandUpC2S>(OnStandUp);
                Handle.UnRegist<MoveC2S>(OnMove);
                Handle.UnRegist<ReadyC2S>(OnReady);
                Handle.UnRegist<CountC2S>(OnCount);
                Handle.UnRegist<GameResultC2S>(OnGameResult);
                Handle.Send(new ExitRoomS2C { Code = m_Room.UserExit(this) });
            }

            private void OnExitRoom(TcpProtoHandle handle, ExitRoomC2S msg) {
                UnRegist();
            }

            private void OnSitDown(TcpProtoHandle handle, SitDownC2S msg) {
                handle.Send(new SitDownS2C { Code = m_Room.SitDown(this, msg.Seat) });
            }

            private void OnStandUp(TcpProtoHandle handle, StandUpC2S msg) {
                handle.Send(new StandUpS2C { Code = m_Room.StandUp(this) });
            }

            private void OnMove(TcpProtoHandle handle, MoveC2S msg) {
                Handle.Send(new MoveS2C { Code = m_Room.Move(this, msg.Seat) });
            }

            private void OnReady(TcpProtoHandle handle, ReadyC2S msg) {
                Handle.Send(new ReadyS2C { Code = m_Room.Ready(this, msg.Value) });
            }

            private void OnCount(TcpProtoHandle handle, CountC2S msg) {
                handle.Send(new CountS2C { Code = m_Room.Count(this, msg.Start) });
            }

            private void OnGameResult(TcpProtoHandle handle, GameResultC2S msg) {
                handle.Send(new GameResultS2C { Code = m_Room.GameResult(this, msg.ScorePBs) });
            }
        }
    }
}
