using BDFramework.UFlux;
using MR.Net.Proto.Battle;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using TMPro;
using UnityEngine;

public class SubRoomWindow : AWindow
{
    public SubRoomWindow(string path) : base(path)
    {

    }

    public SubRoomWindow(Transform transform) : base(transform)
    {

    }

    [TransformPath("RoomIDIpt")]
    public TMP_InputField RoomIDIpt;

    [ButtonOnclick("CloseBtn")]
    public void OnClose()
    {
        this.Close();
    }
    [ButtonOnclick("JoinBtn")]
    public void OnJoinRoom()
    {
        try
        {
            Window_Main.RoomID = int.Parse(RoomIDIpt.text);
            UIManager.Inst.ShowWindow(WinEnum.Win_BattleReady);
            Close();
        }
        catch (Exception e)
        {
            var mData = new Window_Tips.UIMsg_Tips();
            mData.context = e.Message.ToString();
            mData.isCancel = false;
            mData.title = "ERROR";
           // mData.confirm = Close;
            UIManager.Inst.ShowWindow(WinEnum.Win_Tips, mData, true, UILayer.Top);
        }
    }

    public override void Close()
    {
        base.Close();
    }
}

