using MR.Net;
using MR.Net.Proto;
using MR.Net.Proto.Battle;
using System;
using System.Collections.Generic;

namespace MR.BattleServer {
    public class Game {
        private const int GAME_TIME = 180;
        private const int GAME_DELAY_TIME = 30;

        public bool Sync { get; private set; }

        private List<User> m_Players = new List<User>();
        private List<User> m_Watchers = new List<User>();
        private Random m_Ran = new Random();
        private KcpService m_Kcp;

        public State GameState { get; private set; }
        private DateTime m_StartTime;

        private List<byte[]> m_ControlData = new List<byte[]>();
        private int[] playerControls = new int[6];
        private int m_Frame;
        private DateTime m_GameDeadLine;

        private ScoreData[] m_ScoreData = new ScoreData[6];

        private bool m_TestMode;

        public Game(TcpProtoHandle[] playerHandles, List<TcpProtoHandle> watchHandles, int portMin, int portMax) {
            var seats = new List<int>();
            for (int i = 0; i < playerHandles.Length; i++) {
                var handle = playerHandles[i];
                if (handle != null) {
                    m_Players.Add(new User(this, handle, i, m_Ran.Next()));
                    seats.Add(i);
                }
            }
            m_TestMode = !CheckPlayer();
            foreach (var handle in watchHandles)
                m_Watchers.Add(new User(this, handle, -1, m_Ran.Next()));

            int port = portMin;
            while (m_Kcp == null && port <= portMax) {
                try {
                    m_Kcp = new KcpService(port);
                } catch {
                    port++;
                }
            }

            if (m_Kcp == null) {
                //TODO
                Console.WriteLine("Ports occupied.");
                return;
            }

            m_Kcp.onNewClient += OnNewClient;
            var sead = m_Ran.Next();
            var pi = seats.ToArray();
            foreach (var user in m_Players)
                user.Handle.Send(new GameStartS2C { Index = user.Index, Token = user.Token, PlayerIndexs = pi, Sead = sead, Port = port, GameTime = GAME_TIME, GameDelayTime = GAME_DELAY_TIME });
            foreach (var user in m_Watchers)
                user.Handle.Send(new GameStartS2C { Index = user.Index, Token = user.Token, PlayerIndexs = pi, Sead = sead, Port = port, GameTime = GAME_TIME, GameDelayTime = GAME_DELAY_TIME });
        }

        private async void OnNewClient(KcpService.Client client) {
            var res = await client.ReceiveAsync();
            var token = BitConverter.ToInt32(res);
            var target = m_Players.Find(m => m.Token == token);
            if (target == null)
                target = m_Watchers.Find(m => m.Token == token);
            if (target != null)
                target.ReadyKCP(client);
        }

        private bool CheckPlayer() {
            return m_Players.Exists(m => m.Index >= 0 && m.Index < 3) && m_Players.Exists(m => m.Index >= 3 && m.Index < 6);
        }

        public void Update() {
            if (m_Players.Count == 0 || !m_TestMode && !CheckPlayer()) {
                GameState = State.Stop;
                return;
            }
            switch (GameState) {
                case State.Ready:
                    var hasWait = false;
                    lock (m_Players)
                        hasWait = m_Players.Exists(m => m.KcpClient == null) || m_Watchers.Exists(m => m.KcpClient == null);
                    if (!hasWait) {
                        GameState = State.Begin;
                        m_StartTime = DateTime.Now;
                    }
                    break;
                case State.Begin:
                    while ((int)Math.Floor((DateTime.Now - m_StartTime).TotalSeconds * 30) > m_Frame) {
                        lock (m_ControlData) {
                            if (m_ControlData.Count == m_Frame)
                                m_ControlData.Add(BitConverter.GetBytes(m_ControlData.Count));
                            for (int i = 0; i < playerControls.Length; i++)
                                if (playerControls[i] < m_Frame)
                                    playerControls[i] = m_Frame;
                            SendAllKcp(m_ControlData[m_Frame]);
                        }
                        m_Frame++;
                    }
                    break;
                case State.Score:
                    if (DateTime.Now > m_GameDeadLine) {
                        GameState = State.Stop;
                    }
                    break;
            }
        }


        public void Contorl(byte index, byte[] data) {
            //if (data.Length > 1)
            //    Console.WriteLine("Recv:" + string.Join(",", data));
            lock (m_ControlData) {
                if (playerControls[index] >= m_ControlData.Count - 1) {
                    var nData = new byte[5 + data.Length];
                    BitConverter.TryWriteBytes(nData, m_ControlData.Count);
                    nData[4] = index;
                    Array.Copy(data, 0, nData, 5, data.Length);
                    m_ControlData.Add(nData);
                } else {
                    var step = playerControls[index] + 1;
                    var oData = m_ControlData[step];
                    var nData = new byte[oData.Length + 1 + data.Length];
                    Array.Copy(oData, 0, nData, 0, oData.Length);
                    nData[oData.Length] = index;
                    Array.Copy(data, 0, nData, 1 + oData.Length, data.Length);
                    m_ControlData[step] = nData;
                }
            }
            playerControls[index]++;
        }

        public void ReciveScore(TcpProtoHandle handle, ScorePB[] scorePBs) {
            if (GameState == State.Stop)
                return;
            var player = m_Players.Find(m => m.Handle == handle);
            if (player == null || player.Result)
                return;
            player.Result = true;
            if (GameState == State.Score) {
                foreach (var s in scorePBs) {
                    var score = m_ScoreData[s.Index];
                    if (score == null || score.kill != s.Kill || score.die != s.Die || score.output != s.Output || score.take != s.Take) {
                        GameState = State.Stop;
                        return;
                    }
                }
            } else {
                GameState = State.Score;
                foreach (var s in scorePBs)
                    m_ScoreData[s.Index] = new ScoreData { kill = s.Kill, die = s.Die, output = s.Output, take = s.Take };
            }
            if (!m_Players.Exists(m => m.Result)) {
                GameState = State.Stop;
                Sync = true;
            }
        }

        private void SendAllKcp(byte[] data) {
            //if (data.Length > 6)
            //    Console.WriteLine("Send:" + string.Join(",", data));
            foreach (var user in m_Players)
                user.KcpClient.Send(data, data.Length);
            foreach (var user in m_Watchers)
                user.KcpClient.Send(data, data.Length);
        }

        private void Leave(User user) {
            m_Players.Remove(user);
            m_Watchers.Remove(user);
        }

        public void Dispose() {
            m_Kcp.Dispose();
        }

        private class ScoreData {
            public int kill;
            public int die;
            public long output;
            public long take;
        }

        public enum State {
            Ready,
            Begin,
            Score,
            Stop
        }

        public class User {
            public bool Result { get; set; }

            public TcpProtoHandle Handle { get; private set; }
            public int Index { get; private set; }
            public int Token { get; private set; }
            public KcpService.Client KcpClient { get; private set; }

            private Game m_Game;


            public User(Game game, TcpProtoHandle handle, int index, int token) {
                Handle = handle;
                Index = index;
                Token = token;
                m_Game = game;
                Handle.onClose += OnClose;
            }

            public void ReadyKCP(KcpService.Client client) {
                KcpClient = client;
                if (Index != -1)
                    Read();
            }

            private async void Read() {
                while (m_Game.GameState != State.Stop) {
                    if (m_Game.GameState == State.Begin)
                        try {
                            m_Game.Contorl((byte)Index, await KcpClient.ReceiveAsync());
                        } catch { }
                }
            }

            private void OnClose() {
                m_Game.Leave(this);
                Handle.onClose -= OnClose;
            }
        }
    }
}
