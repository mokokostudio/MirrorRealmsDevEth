using System;
using System.Collections.Generic;
using TrueSync;
using TrueSync.Physics3D;
using TSWorld = TrueSync.Physics3D.World;

namespace MR.Battle {
    public class BattleGroundCD : ComponentData {
        public byte[] OperateData { get; set; }
        public OperateDataProcesser OperateDataProcesser { get; set; }
        public Random Random { get; set; }
        public Dictionary<byte, PlayerCD> Players { get; } = new Dictionary<byte, PlayerCD>();
        public PlayerCD PlayerHandle { get; set; }
        public PlayerCD PlayerWatch { get; set; }
        public bool WatchNext { get; set; }
        public TSWorld TSWorld { get; set; }
        public FP PreCameraDir { get; set; }
        public bool Rage { get; set; }
        public MutiMap<RigidBody, RigidBody> collisionMap { get; } = new MutiMap<RigidBody, RigidBody>();
        private Dictionary<Entity, RigidBody> m_E2RDic = new Dictionary<Entity, RigidBody>();
        private Dictionary<RigidBody, Entity> m_R2EDic = new Dictionary<RigidBody, Entity>();
        public void RegistPhysic(Entity e, RigidBody r) {
            TSWorld.AddBody(r);
            m_E2RDic.Add(e, r);
            m_R2EDic.Add(r, e);
        }
        public void UnregisterPhysic(Entity e) {
            var r = m_E2RDic[e];
            TSWorld.RemoveBody(r);
            m_E2RDic.Remove(e);
            m_R2EDic.Remove(r);
        }
        public List<Entity> GetCollisionTarget(Entity e) {
            if (m_E2RDic.TryGetValue(e, out var r))
                return collisionMap.GetValues(r).ConvertAll(m => m_R2EDic[m]);
            else
                return new List<Entity>();
        }
    }
}
