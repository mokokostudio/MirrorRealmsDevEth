using System.Collections;
using System.Collections.Generic;

namespace MR.Battle {
    public class MutiMap<K, V> : IEnumerable<KeyValuePair<K, List<V>>> {
        private readonly Dictionary<K, List<V>> m_Dic = new Dictionary<K, List<V>>();

        public List<V> this[K k] => m_Dic[k];

        public void Add(K k) {
            if (!m_Dic.ContainsKey(k))
                m_Dic.Add(k, new List<V>());
        }

        public void Add(K k, V v) {
            if (!m_Dic.TryGetValue(k, out var list))
                m_Dic[k] = list = new List<V>();
            list.Add(v);
        }

        public void Remove(K k, V v) {
            if (!m_Dic.ContainsKey(k))
                return;
            m_Dic[k].Remove(v);
        }

        public V GetValue(K k) {
            if (TryGetValues(k, out var list) && list.Count > 0)
                return list[0];
            return default;
        }

        public List<V> GetValues(K k) {
            if (TryGetValues(k, out var list))
                return list;
            return new List<V>();
        }

        public bool ContainsKey(K k) => m_Dic.ContainsKey(k);

        public bool TryGetValues(K k, out List<V> list) {
            return m_Dic.TryGetValue(k, out list);
        }

        public IEnumerator<KeyValuePair<K, List<V>>> GetEnumerator() {
            return m_Dic.GetEnumerator();
        }

        IEnumerator IEnumerable.GetEnumerator() {
            return m_Dic.GetEnumerator();
        }

        public void Clear() {
            m_Dic.Clear();
        }
    }
}
