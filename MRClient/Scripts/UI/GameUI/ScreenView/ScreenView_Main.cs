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

[ScreenView((int) ScreenViewEnum.Main)]
public class ScreenView_Main : IScreenView
{
    public int Name { get; set; }
    public bool IsLoad { get; private set; }

    public void BeginInit()
    {
        //一定要设置为true，否则当前是未加载状态
        this.IsLoad = true;
        UIManager.Inst.ShowWindow(WinEnum.Win_Loading);
        UFluxUtils.TaskList.Clear();
        UFluxUtils.UnAsyncLoadScene();

        UIManager.Inst.LoadWindow(WinEnum.Win_Main);
        UIManager.Inst.LoadWindow(WinEnum.Win_BattleReady);
        UIManager.Inst.LoadWindow(WinEnum.Win_Tips);
        // var msg = new UIMsg_Loading();
        // msg.isOpen = true;
        // msg.winEnum = WinEnum.Win_Main;
       
    }

    public void BeginExit()
    {
        UIManager.Inst.UnLoadWindow(WinEnum.Win_Main);
        //UIManager.Inst.UnLoadWindow(WinEnum.Win_BattleReady);
    }
}