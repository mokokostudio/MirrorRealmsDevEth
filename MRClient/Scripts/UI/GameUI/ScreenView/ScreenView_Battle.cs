using System;
using System.Collections;
using System.Collections.Generic;
using BDFramework.ScreenView;
using UnityEngine;
using BDFramework.UFlux;
using UnityEditor;
using UnityEngine.ResourceManagement.AsyncOperations;
using Cysharp.Threading.Tasks;
using MR.Net.Proto.Battle;
using UI.GameUI;
using UnityEngine.AddressableAssets;

[ScreenView((int) ScreenViewEnum.Battle)]
public class ScreenView_Battle : IScreenView
{
    public int Name { get; set; }
    public bool IsLoad { get; private set; }
    private List<Enum> Idxs = new List<Enum>();

    public void BeginInit()
    {
        //一定要设置为true，否则当前是未加载状态
        IsLoad = true;
        UIManager.Inst.ShowWindow(WinEnum.Win_Loading);
        UFluxUtils.TaskList.Clear();
        UIManager.Inst.LoadWindow(WinEnum.Win_Tips);
        UIManager.Inst.LoadWindow(WinEnum.Win_BattleReady);
        var msg = new UIMsg_Loading();
        msg.isOpen = false;
        msg.winEnum = WinEnum.Win_Loading;
        UIManager.Inst.SendMessage(WinEnum.Win_Loading, msg);
        LoadBattle().Forget();
    }

    public async UniTask LoadBattle()
    {
        var handle = Addressables.LoadAssetAsync<GameObject>($"Assets/UI/RoomTestBattle.prefab");
        await handle;
        var uiroot = GameObject.Find("UIRoot")?.transform;
        UnityEngine.GameObject.Instantiate(handle.Result, uiroot.Find("Bottom")?.transform);
    }

    public void BeginExit()
    {
    }
}