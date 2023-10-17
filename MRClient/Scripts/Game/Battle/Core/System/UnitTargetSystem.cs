using TrueSync;

namespace MR.Battle {
    public class UnitTargetSystem : BaseSystem<UnitCD> {
        public override string Group => "Update";
        public override int Order => 800;
        private static FP s_MaxDis = Config.Battle.Constant.LockMissDistance / (FP)100;

        protected override void Run() {
            if (Data.Target == null)
                return;
            if (Data.Target.Entity == null) {
                Data.Target = null;
                return;
            }
            var la = Data.GetComponentData<LocationCD>();
            var lb = Data.Target.GetComponentData<LocationCD>();
            var dis = TSVector.Distance(la.Position, lb.Position);
            if (dis > s_MaxDis)
                Data.Target = null;
        }
    }
}
