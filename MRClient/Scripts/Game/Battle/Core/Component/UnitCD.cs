using System.Collections.Generic;
using TrueSync;
using TrueSync.Physics3D;

namespace MR.Battle {
    public class UnitCD : ComponentData {
        public BattleGroundCD BattleGround { get; set; }
        public PlayerCD Player { get; set; }
        public bool TryMove { get; set; }
        public FP TryFace { get; set; }
        public bool TryDash { get; set; }
        public bool TryAttack { get; set; }
        public bool TryDefense { get; set; }
        public UnitState State { get; set; } = UnitState.Idle;
        public FP MoveSpeed { get; set; }
        public bool LimitSkill { get; set; }
        public bool LimitMove { get; set; }
        public bool LimitDodge { get; set; }
        public bool AttackHolding { get; set; }
        public RigidBody RigidBody { get; set; }
        public LocationCD Target { get; set; }
        public FP HP { get; set; }
        public List<FP> HPChange { get; } = new List<FP>();
        public Camp Camp { get; set; }
    }
}
