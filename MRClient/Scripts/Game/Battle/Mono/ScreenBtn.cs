using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.EventSystems;
using UnityEngine.InputSystem.Layouts;
using UnityEngine.InputSystem.OnScreen;

public class ScreenBtn : OnScreenControl, IDragHandler {
    [InputControl(layout = "Vector2")]
    [SerializeField]
    private string m_ControlPath;

    protected override string controlPathInternal {
        get => m_ControlPath;
        set => m_ControlPath = value;
    }

    public void OnDrag(PointerEventData eventData) {
        SendValueToControl(eventData.delta);
    }
}
