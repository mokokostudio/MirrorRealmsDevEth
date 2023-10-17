namespace MR.Battle {
    public abstract class ComponentData {
        public World World { get; set; }
        public Entity Entity => World.GetEntity(this);
        public T GetComponentData<T>() where T : ComponentData => World.GetComponentData<T>(Entity);
        public T AddComponentData<T>() where T : ComponentData => World.AddComponentData<T>(Entity);
        public T GetOrAddComponentData<T>() where T:ComponentData => World.GetOrAddComponentData<T>(Entity);
        public void RemoveComponentData<T>(T cd) where T : ComponentData => World.RemoveComponentData(Entity, cd);
    }
}
