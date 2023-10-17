using TrueSync;
using TrueSync.Physics3D;

namespace MR.Battle {
    public class AttackBoxCD : ComponentData {
        public BattleGroundCD BattleGround { get; set; }
        public UnitCD Unit { get; set; }
        public TSVector Position { get; set; }
        public TSMatrix Rotation { get; set; }
        public TSVector Scale { get; set; }
        public int AtkLV { get; set; }
        public HitType HitType { get; set; }
        public RigidBody RigidBody { get; set; }
        public int SP { get; set; }
        public int HP { get; set; }
        public int Paush { get; set; }
    }

    public class BulletCD : ComponentData {
        public BattleGroundCD BattleGround { get; set; }
        public UnitCD Unit { get; set; }
        public TSVector Position { get; set; }
        public TSVector Rotation { get; set; }
        public int AtkLV { get; set; }
        public HitType HitType { get; set; }
        public int SP { get; set; }
        public int HP { get; set; }

        public bool FaceTarget { get; set; }
        public FP Speed { get; set; }
        public FP TrackDeg { get; set; }
        public FP TrackTime { get; set; }
    }
}
