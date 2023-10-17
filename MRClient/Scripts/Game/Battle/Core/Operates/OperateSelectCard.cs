using System.IO;

namespace MR.Battle {
    public class OperateSelectCard : IOperateData {
        public byte index;
        public bool Equals(IOperateData other) => other is OperateSelectCard && (other as OperateSelectCard).index == index;
        public void Read(BinaryReader reader) => index = reader.ReadByte();
        public void Write(BinaryWriter writer) => writer.Write(index);
    }
}
