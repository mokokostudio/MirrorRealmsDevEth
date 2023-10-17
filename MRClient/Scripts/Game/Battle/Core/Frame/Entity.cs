namespace MR.Battle {
    public class Entity {
        public World World { get; private set; }

        public Entity(World world) => World = world;

        public T GetComponentData<T>() where T : ComponentData => World.GetComponentData<T>(this);
        public T AddComponentData<T>() where T : ComponentData => World.AddComponentData<T>(this);
        public void RemoveComponentData<T>(T cd) where T : ComponentData => World.RemoveComponentData(this, cd);
    }
}
