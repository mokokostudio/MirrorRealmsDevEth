using System.IO;

namespace MR.Battle {
    public class OperateDefense : IOperateData {
        public bool Equals(IOperateData other) => other is OperateDefense;
        public void Read(BinaryReader reader) { }
        public void Write(BinaryWriter writer) { }
    }
}
