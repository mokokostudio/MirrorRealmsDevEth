using UnityEngine;


    [RequireComponent(typeof(RectTransform))]
    [DisallowMultipleComponent]
    public abstract class AdapterBase : MonoBehaviour
    {
        public abstract void Adapt();
    }
