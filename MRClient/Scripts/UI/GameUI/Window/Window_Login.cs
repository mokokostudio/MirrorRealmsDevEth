
using BDFramework.ScreenView;
using BDFramework.UFlux;
using MR.Net.Frame;
using MR.Net.Proto.Battle;
using System;
using System.Security.Cryptography;
using System.Text;
using TMPro;
using UI.GameUI;
using UnityEngine;

[UI((int)WinEnum.Win_Login, "Window/Window_Login")]
public class Window_Login : AWindow {
    public Window_Login(string path) : base(path) {

    }
    [TransformPath("Adapter/Server/DpdChannel")]
    public TMP_Dropdown DpdChannel;
    [TransformPath("Adapter/Server/AccountID")]
    public TMP_InputField AccountNameInput;
    [TransformPath("Adapter/Server/Password")]
    public TMP_InputField PasswordInput;

    private bool m_Waiting;

    public override void Init() {
        base.Init();
        AccountNameInput.text = PlayerPrefs.GetString("Account");
        for (int i = 0; i < Net.Channels.Count; i++)
            DpdChannel.options.Add(new TMP_Dropdown.OptionData(Net.Channels[i].name));
        DpdChannel.value = PlayerPrefs.GetInt("ChannelIndex");
        DpdChannel.onValueChanged.AddListener(OnChannelChange);
        BGMManager.Play(BGMManager.BGMType.Login);
        UIManager.Inst.LoadWindow(WinEnum.Win_Tips);
    }
    private void OnChannelChange(int index) {
        PlayerPrefs.SetInt("ChannelIndex", index);
    }

    [ButtonOnclick("Adapter/Server/GameButton")]
    private void OnLogin() {
        if (m_Waiting) {
            Debug.Log("Waiting Connect Message.");
            return;
        }
        m_Waiting = true;
        PlayerPrefs.SetString("Account", AccountNameInput.text);
        Net.ChannelIdx = DpdChannel.value;
        try {
            var tcp = Net.CreateTcp();
            tcp.Send<LoginC2S, LoginS2C>(new LoginC2S { Account = AccountNameInput.text, Password = GetMD5(PasswordInput.text) }, OnLoginMsg);
        } catch (Exception e) {
            Debug.LogException(e);
            m_Waiting = false;
        }
    }

    private void OnLoginMsg(LoginS2C msg) {
        m_Waiting = false;
        if (msg.Code == CodePBType.Success) {
            PlayerData.Set(msg.UserInfo);
            ScreenViewManager.Inst.MainLayer.BeginNavTo(ScreenViewEnum.Main);
            UIManager.Inst.SendMessage(WinEnum.Win_Loading, new UIMsg_Loading() { isOpen = true, winEnum = WinEnum.Win_Main });
        } else {
            var mData = new Window_Tips.UIMsg_Tips();
            mData.context = msg.Code.ToString();
            mData.isCancel = false;
            mData.title = "NOTIFY";
            UIManager.Inst.ShowWindow(WinEnum.Win_Tips, mData, true, UILayer.Top);
        }
    }

    private string GetMD5(string input) {
        var md5 = MD5.Create();
        var data = md5.ComputeHash(Encoding.UTF8.GetBytes(input));
        var rule = "";
        for (int i = 0; i < data.Length; i++)
            rule += data[i].ToString("x2");
        return rule;
    }
}