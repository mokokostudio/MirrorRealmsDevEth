using System.IO;

namespace MR.Battle {
    public class OperateAttackEnd : IOperateData {
        public bool Equals(IOperateData other) => other is OperateAttackEnd;

        public void Read(BinaryReader reader) { }
        public void Write(BinaryWriter writer) { }
    }
}
