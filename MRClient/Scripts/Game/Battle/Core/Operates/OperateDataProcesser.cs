using System;
using System.Collections.Generic;
using System.IO;

namespace MR.Battle {
    public class OperateDataProcesser {

        private Dictionary<Type, byte> m_TypeDic = new Dictionary<Type, byte>();
        private byte m_Index = 0;
        private Type[] m_Typs = new Type[256];

        public void Regist<T>() where T : class, IOperateData, new() {
            if (m_Index > 127)
                throw new Exception("Input types is Full");
            AddType(typeof(T), m_Index);
            m_Index++;
        }

        public void RegistTwin<T, U>() where T : class, IOperateData, new() where U : class, IOperateData, new() {
            if (m_Index > 127)
                throw new Exception("Input types is Full");
            AddType(typeof(T), m_Index);
            AddType(typeof(U), (byte)~m_Index);
            m_Index++;
        }

        private void AddType(Type type, byte index) {
            m_TypeDic.Add(type, index);
            m_Typs[index] = type;
        }

        public List<IOperateData> Deserialize(byte[] data, ref int offset) {
            List<IOperateData> inputs = new List<IOperateData>();
            using (var ms = new MemoryStream(data, offset, data.Length - offset)) {
                using (var br = new BinaryReader(ms)) {
                    var num = br.ReadByte();
                    for (byte i = 0; i < num; i++) {
                        var id = br.ReadByte();
                        var type = m_Typs[id];
                        var input = Activator.CreateInstance(type) as IOperateData;
                        input.Read(br);
                        inputs.Add(input);
                    }
                    offset += (int)ms.Position;
                }
            }
            return inputs;
        }

        private List<byte> m_CommandIDs = new List<byte>();
        private List<IOperateData> m_CommandDatas = new List<IOperateData>();

        public void PushInput<T>(T data) where T : class, IOperateData, new() {
            var id = m_TypeDic[typeof(T)];
            var i = m_CommandIDs.IndexOf(id);
            if (i >= 0) {
                m_CommandDatas[i] = data;
                return;
            }
            var uid = (byte)~id;
            var ui = m_CommandIDs.IndexOf(uid);
            if (ui >= 0) {
                m_CommandIDs.RemoveAt(ui);
                m_CommandDatas.RemoveAt(ui);
                return;
            }
            m_CommandIDs.Add(id);
            m_CommandDatas.Add(data);
        }

        public List<IOperateData> TakeByte(out byte[] data) {
            List<IOperateData> inputs = new List<IOperateData>();
            using (var ms = new MemoryStream()) {
                using (var bw = new BinaryWriter(ms)) {
                    bw.Write((byte)m_CommandIDs.Count);
                    for (byte i = 0; i < m_CommandIDs.Count; i++) {
                        var id = m_CommandIDs[i];
                        var input = m_CommandDatas[i];
                        bw.Write(id);
                        input.Write(bw);
                        inputs.Add(input);
                    }
                    m_CommandIDs.Clear();
                    m_CommandDatas.Clear();
                }
                data = ms.ToArray();
            }
            return inputs;
        }
    }
}
