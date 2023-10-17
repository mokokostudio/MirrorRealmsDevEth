using TrueSync;
using TrueSync.Physics3D;

namespace MR.Battle {
    public class UnitPhySystem : BaseSystem<UnitCD> {
        public override string Group => "Update";
        public override int Order => 0;

        protected override void Run() {
            if (Data.RigidBody == null) {
                Data.RigidBody = new RigidBody(new CapsuleShape(1.6, 0.5));
                Data.BattleGround.RegistPhysic(Entity, Data.RigidBody);
            }
            var location = GetComponentData<LocationCD>();
            Data.RigidBody.Position = location.Position;
        }
    }

    public class AttackBoxPhySystem : BaseSystem<AttackBoxCD> {
        public override string Group => "Update";
        public override int Order => 0;

        protected override void Run() {
            if (Data.RigidBody == null) {
                Data.RigidBody = new RigidBody(new BoxShape(Data.Scale));
                Data.BattleGround.RegistPhysic(Entity, Data.RigidBody);
            }
            var location = Data.Unit.GetComponentData<LocationCD>();
            var anim = Data.Unit.GetComponentData<UnitAnimCD>();
            var face = (anim.LookTarget ? anim.Face : location.Face).AsFloat();
            var mat = TSMatrix.CreateRotationY(face);
            Data.RigidBody.Position = TSVector.Transform(Data.Position, mat) + location.Position;
            Data.RigidBody.Orientation = mat * Data.Rotation;
        }
    }
}
