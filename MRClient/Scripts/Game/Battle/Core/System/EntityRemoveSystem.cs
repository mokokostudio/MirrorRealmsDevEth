using Sirenix.Utilities;

namespace MR.Battle {

    public class EntityRemoveSystem : BaseSystem<EntityRemoveCD> {
        public override string Group => "Update";
        public override int Order => 1100;

        protected override void Run() {
            World.RemoveEntity(Entity);
        }
    }
}
