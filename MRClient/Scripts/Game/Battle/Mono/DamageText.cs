using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;

public class DamageText : MonoBehaviour {
    public int Damage { get; set; }
    public int index { get; set; }
    public TMP_Text damageTxt;

    private void Start() {
        if (Damage != 0) {
            damageTxt.transform.localPosition += index * new Vector3(0, 30, 30);
            damageTxt.text = Damage.ToString();
            damageTxt.color = Damage > 0 ? Color.green : Color.red;
        }
        Destroy(gameObject, 1);
    }

    private void Update() {
        transform.rotation = Camera.main.transform.rotation;
    }
}
