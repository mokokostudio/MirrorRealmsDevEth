using System.IO;

namespace MR.Battle {
    public class OperateDiscardWeapon : IOperateData {
        public bool Equals(IOperateData other) => other is OperateDiscardWeapon;
        public void Read(BinaryReader reader) { }
        public void Write(BinaryWriter writer) { }
    }
}
