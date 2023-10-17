using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class FreeItem : MonoBehaviour {
    public static void Free(GameObject go, float time) {
        Destroy(go, time);
        go.AddComponent<FreeItem>();
    }

    private void Start() {
        transform.SetParent(null);
        gameObject.AddComponent<Rigidbody>();
    }
}
