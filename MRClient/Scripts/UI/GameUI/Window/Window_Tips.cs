
using BDFramework.ScreenView;
using BDFramework.UFlux;
using MR.Net.Frame;
using MR.Net.Proto.Battle;
using System;
using TMPro;
using UnityEngine;
using UnityEngine.UI;

[UI((int)WinEnum.Win_Tips, "Window/Window_Tips")]
public class Window_Tips : AWindow
{
    public Window_Tips(string path) : base(path)
    {

    }

    public override void Init()
    {
        base.Init();

    }
    public override void Open(UIMsgData uiMsg = null)
    {
        base.Open(uiMsg);
        Message(uiMsg as UIMsg_Tips);
    }
    [TransformPath("CancelButton")]
    public Transform CancelButton;
    [TransformPath("ConfirmButton")]
    public Transform ConfirmButton;
    [TransformPath("CConfirmButton")]
    public Transform CConfirmButton;
    [TransformPath("CConfirmButton/CloseBtn")]
    public Transform CloseBtn;
    [TransformPath("Context")]
    public TMP_Text Context;
    [TransformPath("Title")]
    public TMP_Text Title;

    public class UIMsg_Tips : UIMsgData
    {
        public bool isCancel;
        public string context;
        public string title;
        public Action cancel;
        public Action confirm;
        public Action callBack;
    }
    [UIMessageListener]
    private void Message(UIMsg_Tips msg)
    {
        Context.text = msg.context;
        Title.text = msg.title;
        CConfirmButton.GetComponent<Button>().onClick.RemoveAllListeners();
        ConfirmButton.GetComponent<Button>().onClick.RemoveAllListeners();
        CancelButton.GetComponent<Button>().onClick.RemoveAllListeners();
        CConfirmButton.GetComponent<Button>().onClick.AddListener(() =>
        {
            msg.confirm?.Invoke();
            msg.callBack?.Invoke();
            Close();
        });
        ConfirmButton.GetComponent<Button>().onClick.AddListener(() =>
        {
            msg.confirm?.Invoke();
            msg.callBack?.Invoke();
            Close();
        });
        CancelButton.GetComponent<Button>().onClick.AddListener(() =>
        {
            msg.cancel?.Invoke();
            msg.callBack?.Invoke();
            Close();
        });
        CloseBtn.GetComponent<Button>().onClick.AddListener(() =>
        {
            Close();
        });
        if (msg.isCancel)
        {
            CConfirmButton.gameObject.SetActive(false);
            ConfirmButton.gameObject.SetActive(true);
            CancelButton.gameObject.SetActive(true);
        }
        else
        {
            CConfirmButton.gameObject.SetActive(true);
            ConfirmButton.gameObject.SetActive(false);
            CancelButton.gameObject.SetActive(false);
        }
    }
}