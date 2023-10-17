using System;

namespace MR.Battle {
    public abstract class BaseSystem<T> : ISystem where T : ComponentData {
        public World World { get; set; }
        public Type Type => typeof(T);
        public abstract string Group { get; }
        public abstract int Order { get; }

        public void Run(object data) {
            Data = (T)data;
            Run();
        }
        protected abstract void Run();

        protected T Data { get; private set; }
        protected Entity Entity => Data.Entity;

        protected U GetComponentData<U>() where U : ComponentData => Data.GetComponentData<U>();
        protected U AddComponentData<U>() where U : ComponentData => Data.AddComponentData<U>();
        protected void RemoveComponentData<U>(U cd) where U : ComponentData => Data.RemoveComponentData(cd);
    }
}
