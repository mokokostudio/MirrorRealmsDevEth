using System;
using System.Collections.Generic;
using System.Reflection;

namespace MR.Battle {
    public class World {
        private readonly MutiMap<string, ISystem> m_Systems = new MutiMap<string, ISystem>();

        private readonly Dictionary<Entity, MutiMap<Type, ComponentData>> m_EntityComponentMap = new Dictionary<Entity, MutiMap<Type, ComponentData>>();
        private readonly Dictionary<ComponentData, Entity> m_ComponentDataEntityDic = new Dictionary<ComponentData, Entity>();
        private readonly MutiMap<Type, ComponentData> m_ComponentMap = new MutiMap<Type, ComponentData>();

        private readonly List<ComponentData> m_TempComponentDatas = new List<ComponentData>();

        private readonly MutiMap<Type, Type> m_ComponentDataAttachs = new MutiMap<Type, Type>();

        public void LoadAssembly(Assembly assembly) {
            var sysType = typeof(ISystem);
            var cdType = typeof(ComponentData);
            foreach (var type in assembly.GetExportedTypes()) {
                if (type.IsClass && !type.IsAbstract && sysType.IsAssignableFrom(type)) {
                    var sys = (ISystem)Activator.CreateInstance(type);
                    sys.World = this;
                    if (!m_Systems.ContainsKey(sys.Group))
                        m_Systems.Add(sys.Group);
                    m_Systems[sys.Group].Add(sys);
                }

                if (type.IsClass && !type.IsAbstract && cdType.IsAssignableFrom(type)) {
                    foreach (var att in type.GetCustomAttributes<AttachToAttribute>()) {
                        var target = att.Target;
                        if (cdType.IsAssignableFrom(target))
                            m_ComponentDataAttachs.Add(target, type);
                    }
                }
            }

            foreach (var kv in m_Systems)
                kv.Value.Sort((m, n) => n.Order.CompareTo(m.Order));

#if UNITY_STANDALONE
            foreach (var kv in m_Systems)
                foreach (var sys in kv.Value)
                    UnityEngine.Debug.Log($"Group:{kv.Key}, Order:{sys.Order}, System:{sys}");
#endif
        }

        public void Run(string group) {
            if (!m_Systems.TryGetValues(group, out var list))
                return;
            foreach (var sys in list) {
                var type = sys.Type;
                if (m_ComponentMap.TryGetValues(type, out var cds)) {
                    m_TempComponentDatas.AddRange(cds);
                    foreach (var cd in m_TempComponentDatas)
                        sys.Run(cd);
                    m_TempComponentDatas.Clear();
                }
            }
        }

        public Entity CreateEntity() {
            var entity = new Entity(this);
            m_EntityComponentMap.Add(entity, new MutiMap<Type, ComponentData>());
            return entity;
        }

        public void RemoveEntity(Entity entity) {
            m_EntityComponentMap.TryGetValue(entity, out var map);
            foreach (var kv in map) {
                foreach (var cd in kv.Value) {
                    m_ComponentMap.Remove(cd.GetType(), cd);
                    m_ComponentDataEntityDic.Remove(cd);
                }
            }
            m_EntityComponentMap.Remove(entity);
        }

        internal T AddComponentData<T>(Entity e) where T : ComponentData {
            if (!m_EntityComponentMap.TryGetValue(e, out var map))
                throw new Exception("Entity is not exist when use AddComponentData.");
            var type = typeof(T);
            var cd = AddComponentData(e, map, type);
            if (m_ComponentDataAttachs.TryGetValues(type, out var list))
                foreach (var follow in list)
                    AddComponentData(e, map, follow);
            return cd as T;
        }

        private ComponentData AddComponentData(Entity e, MutiMap<Type, ComponentData> map, Type type) {
            var cd = Take(type);
            map.Add(type, cd);
            m_ComponentMap.Add(type, cd);
            m_ComponentDataEntityDic.Add(cd, e);
            return cd;
        }

        internal T GetComponentData<T>(Entity e) where T : ComponentData {
            if (!m_EntityComponentMap.TryGetValue(e, out var map))
                throw new Exception("Entity is not exist when use GetComponentData.");
            var type = typeof(T);
            return (T)map.GetValue(type);
        }

        internal T GetOrAddComponentData<T>(Entity e) where T : ComponentData {
            var cd = GetComponentData<T>(e);
            if (cd == null)
                return AddComponentData<T>(e);
            return cd;
        }

        internal void RemoveComponentData<T>(Entity e, T cd) where T : ComponentData {
            if (!m_EntityComponentMap.TryGetValue(e, out var map))
                throw new Exception("Entity is not exist when use RemoveComponentData.");
            RemoveComponentData(e, map, cd);
        }

        private void RemoveComponentData(Entity e, MutiMap<Type, ComponentData> map, ComponentData cd) {
            var type = cd.GetType();
            map.Remove(type, cd);
            m_ComponentMap.Remove(type, cd);
            m_ComponentDataEntityDic.Remove(cd);
        }

        private ComponentData Take(Type type) {
            var cd = Activator.CreateInstance(type) as ComponentData;
            cd.World = this;
            return cd;
        }

        internal Entity GetEntity(ComponentData cd) {
            m_ComponentDataEntityDic.TryGetValue(cd, out var e);
            return e;
        }

        public List<Entity> GetEntities<T>() where T : ComponentData {
            var type = typeof(T);
            if (m_ComponentMap.TryGetValues(type, out var list))
                return list.ConvertAll(x => x.Entity);
            return new List<Entity>();
        }
    }
}
