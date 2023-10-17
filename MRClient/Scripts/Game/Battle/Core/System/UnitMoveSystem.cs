using TrueSync;

namespace MR.Battle {
    public class UnitMoveSystem : BaseSystem<UnitCD> {
        public override string Group => "Update";
        public override int Order => 500;
        public FP Interval => Battle.Instance.Interval;
        private static FP s_DashRate = (FP)1.5;
        protected override void Run() {
            var offset = TSVector.zero;
            var location = Data.GetComponentData<LocationCD>();
            switch (Data.State) {
                case UnitState.Run:
                    offset = new TSVector(0, 0, Data.MoveSpeed * Interval);
                    location.Face = Data.TryFace;
                    break;
                case UnitState.Dash:
                    offset = new TSVector(0, 0, Data.MoveSpeed * s_DashRate * Interval);
                    location.Face = Data.TryFace;
                    break;
                case UnitState.Die:
                case UnitState.Skill:
                    offset = Data.GetComponentData<UnitAnimCD>().AnimOffset;
                    break;
            }
            if (offset != TSVector.zero) {
                var move = Data.GetComponentData<MoveCD>();
                var x = offset.z * FP.Sin(location.Face) + offset.x * FP.Cos(location.Face);
                var z = offset.z * FP.Cos(location.Face) - offset.x * FP.Sin(location.Face);
                move.Offset = new TSVector(x, 0, z);
            }
        }
    }
}
