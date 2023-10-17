using MR.Net.Proto.Battle;
using System;
using System.Collections.Generic;
using System.Reflection;
using TrueSync;

namespace MR.Battle {
    public class Battle {
        public static Battle Instance { get; private set; }

        public OperateDataProcesser OperateDataProcesser => m_BattleGround.OperateDataProcesser;
        public bool GameOver { get; private set; }
        public FP RemainTime => Math.Max(m_GameFrame - Frame, 0) * Interval;

        private World m_World;
        private readonly Queue<byte[]> m_FrameDatas = new Queue<byte[]>();
        private FP m_Buffer;
        public FP Interval { get; private set; }
        public int Frame { get; private set; }
        public int DataFrame => m_FrameDatas.Count;
        public int ActFrame => -TSMath.Ceiling(m_LogicTime / Interval).AsInt();
        private FP m_LogicTime;
        private FP m_OperateTime;

        private int m_GameFrame;
        private int m_OverFrame;
        private int m_OverTimes;


        private BattleGroundCD m_BattleGround;

        public bool Over => Frame >= m_GameFrame;
        public bool NearOver => Frame + m_OverFrame >= m_GameFrame;

        public BattleGroundScoreCD Score => m_BattleGround.GetComponentData<BattleGroundScoreCD>();

        public FP PreCameraDir {
            get => m_BattleGround.PreCameraDir;
            set => m_BattleGround.PreCameraDir = value;
        }
        public FP CameraDir {
            get {
                if (m_BattleGround.PlayerWatch != null)
                    return m_BattleGround.PlayerWatch.CamDir;
                else
                    return m_BattleGround.PreCameraDir;
            }
        }

        public PlayerCD CameraPlayer {
            get {
                if (m_BattleGround.PlayerWatch != null)
                    return m_BattleGround.PlayerWatch;
                else
                    return m_BattleGround.PlayerHandle;
            }
        }

        public Battle(int frameRate, int buffer, int gameTime, int overTime, int overTimes, params Assembly[] assemblys) {
            Instance = this;
            m_GameFrame = frameRate * gameTime;
            m_OverFrame = frameRate * overTime;
            m_OverTimes = overTimes;
            Interval = 1 / (FP)frameRate;
            m_Buffer = buffer;
            m_World = new World();
            foreach (var item in assemblys)
                m_World.LoadAssembly(item);
        }

        public void Init(GameStartS2C s2c, TSVector[] bornPoints, Dictionary<int, ulong[]> weaponDic) {
            var entity = m_World.CreateEntity();
            var init = entity.AddComponentData<BattleGroundInitCD>();
            init.PlayerIndex = s2c.PlayerIndexs;
            init.HandleIndex = s2c.Index;
            init.Sead = s2c.Sead;
            init.BornPoints = bornPoints;
            init.WeaponDic = weaponDic;
            m_BattleGround = entity.AddComponentData<BattleGroundCD>();
        }

        public List<byte[]> Update(FP deltaTime) {
            var result = new List<byte[]>();

            var rFrame = (FP)m_FrameDatas.Count;
            if (rFrame > m_Buffer)
                m_LogicTime -= TSMath.Lerp(deltaTime, Interval, (rFrame - m_Buffer - 1) / 10);
            else
                m_LogicTime -= deltaTime * rFrame / m_Buffer;
            if (m_LogicTime <= 0) {
                while (m_LogicTime <= 0 && m_FrameDatas.TryDequeue(out var fData)) {
                    Frame++;
                    m_LogicTime += Interval;
                    m_BattleGround.OperateData = fData;
                    m_World.Run("Update");
                }
            }
            m_World.Run("Display");
            m_OperateTime -= deltaTime;
            while (m_OperateTime <= 0 && OperateDataProcesser != null) {
                OperateDataProcesser.TakeByte(out var data);
                m_OperateTime += Interval;
                if (data != null && data.Length > 1)
                    result.Add(data);
            }
            m_BattleGround.Rage = NearOver;
            if (Over) {
                if (Score.killA == Score.killB) {
                    if (m_OverTimes > 0) {
                        m_OverTimes--;
                        m_GameFrame += m_OverFrame;
                    }
                } else {
                    GameOver = true;
                }
            }
            return result;
        }

        public void OperateInput(byte[] datas) {
            m_FrameDatas.Enqueue(datas);
        }

        public void WatchNext() {
            m_BattleGround.WatchNext = true;
        }
    }
}
