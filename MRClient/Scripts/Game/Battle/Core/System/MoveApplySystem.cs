using TrueSync;

namespace MR.Battle {
    public class MoveApplySystem : BaseSystem<MoveCD> {
        public override string Group => "Update";
        public override int Order => 100;
        private static FP s_Side = 24.5;

        protected override void Run() {
            var unit = Data.GetComponentData<UnitCD>();
            var location = Data.GetComponentData<LocationCD>();
            var tarPos = location.Position += Data.Offset;
            var entities = World.GetEntities<UnitCD>();

            if (unit.State != UnitState.Die) {
                var pushv = TSVector.zero;
                foreach (var entity in entities) {
                    var l = entity.GetComponentData<LocationCD>();
                    if (l == location)
                        continue;
                    var delta = location.Position - l.Position;
                    var dis = delta.magnitude;
                    if (dis < 1) {
                        var v = delta.normalized;
                        pushv += v * (1 - dis) / 2;
                    }
                }
                tarPos += pushv;
            }

            tarPos.x = TSMath.Clamp(tarPos.x, -s_Side, s_Side);
            tarPos.z = TSMath.Clamp(tarPos.z, -s_Side, s_Side);
            location.Position = tarPos;
            Data.Offset = TSVector.zero;
        }
    }
}
