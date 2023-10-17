using System;

namespace MR.Battle {
    public class AttachToAttribute : Attribute {
        public Type Target { get; private set; }
        public AttachToAttribute(Type type) {
            Target = type;
        }
    }
}
