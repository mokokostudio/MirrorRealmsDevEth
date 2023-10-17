using System.IO;

namespace MR.Battle {
    public class OperateAttackStart : IOperateData {

        public bool Equals(IOperateData other) => other is OperateAttackStart;

        public void Read(BinaryReader reader) { }
        public void Write(BinaryWriter writer) { }
    }
}
