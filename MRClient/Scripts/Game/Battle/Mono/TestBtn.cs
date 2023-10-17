using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.EventSystems;
using UnityEngine.InputSystem.Layouts;
using UnityEngine.InputSystem.OnScreen;

public class TestBtn : OnScreenControl, IPointerEnterHandler, IPointerExitHandler, IPointerDownHandler, IPointerUpHandler {
    private static TestBtn s_CurrentDown;
    private static TestBtn s_CurrentTarget;

    [InputControl(layout = "Button")]
    [SerializeField]
    private string m_ControlPath;

    private bool m_Down;

    protected override string controlPathInternal {
        get => m_ControlPath;
        set => m_ControlPath = value;
    }

    public void OnPointerEnter(PointerEventData eventData) {
        if (s_CurrentDown != null && s_CurrentDown != this) {
            s_CurrentDown.Release();
            SendValueToControl(1.0f);
            s_CurrentTarget = this;
        } else if (s_CurrentDown == this && !m_Down) {
            s_CurrentTarget = null;
            m_Down = true;
            SendValueToControl(1.0f);
        }
    }

    public void OnPointerExit(PointerEventData eventData) {
        if (s_CurrentDown != null && s_CurrentDown != this) {
            s_CurrentDown.Release();
            SendValueToControl(0.0f);
            s_CurrentTarget = null;
        }
    }

    public void OnPointerUp(PointerEventData eventData) {
        if (m_Down)
            SendValueToControl(0.0f);
        else 
            s_CurrentTarget?.SendValueToControl(0.0f);
        s_CurrentTarget = null;
        s_CurrentDown = null;
        m_Down = false;
    }

    public void OnPointerDown(PointerEventData eventData) {
        SendValueToControl(1.0f);
        s_CurrentDown = this;
        m_Down = true;
    }

    public void Release() {
        if (m_Down) {
            m_Down = false;
            SendValueToControl(0.0f);
        }
    }
}
