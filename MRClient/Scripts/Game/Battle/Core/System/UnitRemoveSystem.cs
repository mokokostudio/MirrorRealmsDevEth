namespace MR.Battle {
    public class UnitRemoveSystem : BaseSystem<EntityRemoveCD> {
        public override string Group => "Update";
        public override int Order => 1200;

        protected override void Run() {
            var unit = GetComponentData<UnitCD>();
            if (unit == null)
                return;
            unit.BattleGround.UnregisterPhysic(Entity);
            foreach (var e in GetComponentData<UnitAnimCD>().AttackBoxEntities.Values)
                unit.BattleGround.UnregisterPhysic(e);
        }
    }
}
