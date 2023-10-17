using System;
using System.Collections;
using System.Collections.Generic;
using BDFramework.ScreenView;
using UnityEngine;
using BDFramework.UFlux;
using UnityEditor;
using UnityEngine.ResourceManagement.AsyncOperations;
using Cysharp.Threading.Tasks;
using UI.GameUI;
using UnityEngine.AddressableAssets;

[ScreenView((int)ScreenViewEnum.Login)]
public class ScreenView_Login : IScreenView
{
    public int Name { get; set; }
    public bool IsLoad { get; private set; }
    private List<Enum> Idxs = new List<Enum>();
    public void BeginInit()
    {
        Application.targetFrameRate = 120;
        Application.runInBackground = true;
        //一定要设置为true，否则当前是未加载状态
        this.IsLoad = true;
        LoadConfig();
        UIManager.Inst.LoadWindow(WinEnum.Win_Loading, () =>
        {
            UFluxUtils.TaskList.Clear();
            Idxs.Add(WinEnum.Win_Login);
            UIManager.Inst.LoadWindow(WinEnum.Win_Login);
            var msg = new UIMsg_Loading();
            msg.isOpen = true;
            msg.winEnum = WinEnum.Win_Login;
            UIManager.Inst.SendMessage(WinEnum.Win_Loading, msg);
        });
       
    }
    public void LoadConfig()
    {
        Config.LoadData(m =>
        {
            var ao = Addressables.LoadAssetAsync<TextAsset>($"Assets/Config/Groups/{m}.bytes");
            ao.WaitForCompletion();
            return ao.Result.bytes;
        });
        var ao2 = Addressables.LoadAssetAsync<TextAsset>($"Assets/Config/Languages/Cn.bytes");
        ao2.WaitForCompletion();
        Config.LoadLanguage(ao2.Result.bytes);
    }
    public void BeginExit()
    {
        UIManager.Inst.CloseWindow(WinEnum.Win_Login);
        UIManager.Inst.UnLoadWindows(Idxs);
    }
}