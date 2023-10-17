using BDFramework.UFlux;
using System;
using System.Collections.Generic;
using System.Threading;
using MR.Net.Frame;
using MR.Net.Proto.Battle;
using UnityEngine;
using UnityEngine.UI;
using BDFramework.ScreenView;
using Cysharp.Threading.Tasks;
using UI.GameUI;
using TMPro;
using BDFramework.UFlux.View.Props;
using BDFramework.UFlux.Collections;
using System.Linq;

public class Room_Item : APropsBase
{
    [ComponentValueBind("Gray", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool playerGroup;
    [ComponentValueBind("ArrowBtn", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool ArrowGo;

    [ComponentValueBind("Gray/Name", typeof(TMP_Text), nameof(TMP_Text.text))]
    public string txtName;
    [ComponentValueBind("Gray/Name", typeof(TMP_Text), nameof(TMP_Text.color))]
    public Color txtNameColor;

    [ComponentValueBind("Gray/ReplaceBtn", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool ReplaceBtn;
    [ComponentValueBind("Gray/ReplaceBtn", typeof(Button), nameof(Button.interactable))]
    public bool reInteractable;

    [ComponentValueBind("Gray/Name/Homeowner", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool Homeowner;

    [ComponentValueBind("Gray/Select", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool Select;

    [ComponentValueBind("ArrowBtn", typeof(Button), nameof(Button.onClick))]
    public Action arrowAction;
    [ComponentValueBind("ArrowBtn", typeof(Button), nameof(Button.interactable))]
    public bool arInteractable;
    [ComponentValueBind("ArrowBtn/Disable", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool aBtnDisable;
    //[ComponentValueBind("Gray/Ready", typeof(TransformHelper),
    //    nameof(TransformHelper.ShowHideTransForm))]
    //public bool readyGo;

    [ComponentValueBind("Gray/HeadIcon", typeof(Image), nameof(Image.sprite))]
    public string headIcon;

    [ComponentValueBind("Gray/Ready", typeof(TransformHelper),
        nameof(TransformHelper.ShowHideTransForm))]
    public bool showReady;

    [ComponentValueBind("Gray/WeaponBtn1", typeof(Button), nameof(Button.onClick))]
    public Action showWeaponAction;
    [ComponentValueBind("Gray/WeaponBtn2", typeof(Button), nameof(Button.onClick))]
    public Action showWeaponAction2;
    [ComponentValueBind("Gray/WeaponBtn3", typeof(Button), nameof(Button.onClick))]
    public Action showWeaponAction3;
    [ComponentValueBind("Gray/WeaponBtn1/Icon", typeof(Image), nameof(Image.sprite))]
    public string weapon1Icon;
    [ComponentValueBind("Gray/WeaponBtn2/Icon", typeof(Image), nameof(Image.sprite))]
    public string weapon2Icon;
    [ComponentValueBind("Gray/WeaponBtn3/Icon", typeof(Image), nameof(Image.sprite))]
    public string weapon3Icon;


    [ComponentValueBind("Gray/WeaponBtn1/nft", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool weapon1Nft;
    [ComponentValueBind("Gray/WeaponBtn2/nft", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool weapon2Nft;
    [ComponentValueBind("Gray/WeaponBtn3/nft", typeof(TransformHelper), nameof(TransformHelper.ShowHideTransForm))]
    public bool weapon3Nft;
}

public class Join_Item : APropsBase
{
    [ComponentValueBind("Name", typeof(TMP_Text), nameof(TMP_Text.text))]
    public string GameNameId;
    [ComponentValueBind("Name", typeof(TMP_Text), nameof(TMP_Text.color))]
    public Color GameNameColor;
    [ComponentValueBind("Name/Owner", typeof(TransformHelper),
        nameof(TransformHelper.ShowHideTransForm))]
    public bool Owner;
}

public class Player
{
    public string name;
    public int seat;
    public bool ready;
    public ulong Character;
    public ulong[] Weapons;
}

public class Room_Window : APropsBase
{
    [ComponentValueBind("Adapter/TabGroup", typeof(UFluxBindLogic), nameof(UFluxBindLogic.BindChildren))]
    public PropsList<Room_Item> ReadyList = new PropsList<Room_Item>();

    [ComponentValueBind("Adapter/OBGroup/ScrollRect/Viewport/Content", typeof(UFluxBindLogic),
        nameof(UFluxBindLogic.BindChildren))]
    public PropsList<Join_Item> JoinList = new PropsList<Join_Item>();
}
[UI((int)WinEnum.Win_BattleReady, "Window/Window_Ready")]
public class Window_BattleReady : AWindow<Room_Window>
{
    public Window_BattleReady(string path) : base(path)
    {
    }

    public TcpInstance TcpInstance;
    public Player Self;
    public List<Player> Players = new List<Player>();
    public int ReadyNum;
    public float WaitTime;
    public bool Gaming;
    private bool isOpen = false;
    private bool isOwner = false;
    private bool isReady = false;
    [TransformPath("Adapter/ReadyButton")]
    public Transform ReadyG;
    [TransformPath("Adapter/StartButton")]
    public Transform StartG;
    [TransformPath("Adapter/UNReadyButton")]
    public Transform UnReadyG;
    [TransformPath("Adapter/CancelButton")]
    public Transform CancelG;
    [TransformPath("Adapter/OBGroup/JoinButton")]
    public Transform JoinG;
    [TransformPath("Adapter/OBGroup/JoinButton/Disable")]
    public Transform JoinGDisable;
    [TransformPath("Adapter/Tips")]
    public Transform Tips;
    [TransformPath("Adapter/Tips/Context")]
    public TMP_Text Tip_Context;

    [TransformPath("Adapter/CutDownImg")] public Image CutDownImg;

    [TransformPath("Adapter/RoomName/RoomID")] public TMP_Text RoomTxt;

    public override void Open(UIMsgData uiMsg = null)
    {
        if (isOpen)
            return;
        TcpInstance = Net.GetTcp();
        TcpInstance.Send(new EnterRoomC2S { RommID = Window_Main.RoomID }, (EnterRoomS2C msg) => OnEnterRoom(msg));
        ReadyG.gameObject.SetActive(false);
        UnReadyG.gameObject.SetActive(false);
        StartG.gameObject.SetActive(false);
        CancelG.gameObject.SetActive(false);
        CutDownImg.gameObject.SetActive(false);
    }

    private void OnEnterRoom(EnterRoomS2C msg)
    {
        if (msg.Code != CodePBType.Success)
        {
            var mData = new Window_Tips.UIMsg_Tips();
            mData.context = msg.Code.ToString();
            mData.isCancel = false;
            mData.confirm = Close;
            mData.title = "ERROR";
            UIManager.Inst.ShowWindow(WinEnum.Win_Tips, mData, true, UILayer.Top);
            Debug.Log(msg.Code);
            TcpInstance = null;
        }
        else
        {
            if (msg.RoomID != 0)
                RoomTxt.text = msg.RoomID.ToString();
            else
                RoomTxt.text = Window_Main.RoomID.ToString();
            for (int i = 0; i < 6; i++)
            {
                var index = i;
                var item = new Room_Item();
                item.playerGroup = false;
                item.ArrowGo = true;
                item.ReplaceBtn = false;
                item.Select = false;
                item.Homeowner = false;
                item.arrowAction = () => Arrow(index);
                //item.readyGo = false;
                item.txtName = "";
                Props.ReadyList.Add(item);
                var joinItem = new Join_Item();
                joinItem.GameNameId = "";
                joinItem.GameNameColor = Color.white;
                Props.JoinList.Add(joinItem);
            }
            Regist();
            Props.SetPropertyChange(nameof(Props.JoinList));
            Props.SetPropertyChange(nameof(Props.ReadyList));
            CommitProps();
        }
    }

    public override void Init()
    {
        base.Init();
      //  RegisterSubWindow(new SubWeaponWindow(this.Transform.Find("Adapter/SubWeaponPanel")));
    }

    private int lastLoadedIndex;
    private CancellationTokenSource cancellationTokenSource;

    private async UniTaskVoid ReadyAll()
    {
        cancellationTokenSource = new CancellationTokenSource();
        while (!cancellationTokenSource.IsCancellationRequested)
        {
            WaitTime -= Time.deltaTime;
            if (WaitTime > 0)
            {
                var index = Mathf.CeilToInt(WaitTime);
                if (index != lastLoadedIndex)
                {
                    lastLoadedIndex = index;
                    try
                    {
                        var sprite =
                            await UFluxUtils.LoadAssetUniTask<Sprite>($"Assets/Arts/UI/Number/Number_{index}.png");
                        CutDownImg.sprite = sprite;
                        CutDownImg.gameObject.SetActive(true);
                    }
                    catch (OperationCanceledException)
                    {
                        // 捕获取消请求异常，终止任务的执行
                        return;
                    }
                }
            }

            await UniTask.Yield();
        }
    }

    private void InitPlayer(PlayerPB[] players)
    {
        Players.Clear();
        ulong index = 0;
        foreach (var data in players)
        {

            var player = new Player
            {
                name = data.Name,
                seat = data.Seat,
                ready = data.Ready,
                Character = index,
                Weapons = data.Weapons
            };
            Players.Add(player);
            if (PlayerData.Name == player.name)
            {
                Self = player;
                isOwner = index == 0;
            }
            index++;
        }

        UpdateSeats();
        base.Open();
        isOpen = true;
        UIManager.Inst.CloseWindow(WinEnum.Win_Main);
    }

    private void UpdateSeats()
    {
        for (int i = 0; i < Props.ReadyList.Count; i++)
        {
            var itemRoom = Props.ReadyList.Get(i);
            var player = Players.Find(m => m.seat == i);
            itemRoom.txtName = player?.name;
            // itemRoom.playerGroup = Players.Find(m => m.seat != -1) != null;
            if (player != null && player.seat != -1)
            {
                itemRoom.txtNameColor = player.Character == 0 ? new Color(0, 1, 1) : Color.white; ;
                itemRoom.ReplaceBtn = Self == player;

                itemRoom.playerGroup = true;
                itemRoom.ArrowGo = false;
                itemRoom.headIcon = "Assets/Arts/UI/BattleReady/head";
                itemRoom.showReady = player.ready || player.Character == 0;
                itemRoom.Homeowner = player.Character == 0;//房主图标
                itemRoom.Select = Self == player;
                if (Self == player)
                {
                    if (player.ready)
                    {
                        itemRoom.reInteractable = false;
                        itemRoom.arInteractable = false;
                        itemRoom.aBtnDisable = false;
                        JoinG.GetComponent<Button>().interactable = false;
                        JoinGDisable.gameObject.SetActive(true);
                    }
                    else
                    {
                        itemRoom.reInteractable = true;
                        itemRoom.arInteractable = true;
                        itemRoom.aBtnDisable = true;
                        JoinG.GetComponent<Button>().interactable = true;
                        JoinGDisable.gameObject.SetActive(false);
                        //itemRoom.showWeaponAction = () => { OnWeaponWindow(0); };
                        //itemRoom.showWeaponAction2 = () => { OnWeaponWindow(1); };
                        //itemRoom.showWeaponAction3 = () => { OnWeaponWindow(2); };
                    }
                }
                itemRoom.weapon1Icon = ResourcePath.WeaponPath + Config.Equips.Weapon[player.Weapons[0]].Icon;
                itemRoom.weapon2Icon = ResourcePath.WeaponPath + Config.Equips.Weapon[player.Weapons[1]].Icon;
                itemRoom.weapon3Icon = ResourcePath.WeaponPath + Config.Equips.Weapon[player.Weapons[2]].Icon;

                itemRoom.weapon1Nft = Config.Equips.Weapon[player.Weapons[0]].Quality > Config.Global.QualityType.N;
                itemRoom.weapon2Nft = Config.Equips.Weapon[player.Weapons[1]].Quality > Config.Global.QualityType.N;
                itemRoom.weapon3Nft = Config.Equips.Weapon[player.Weapons[2]].Quality > Config.Global.QualityType.N;
            }
            else
            {
                itemRoom.showReady = false;
                itemRoom.playerGroup = false;
                itemRoom.ArrowGo = true;
                itemRoom.arInteractable = true;
                itemRoom.aBtnDisable = true;
                itemRoom.txtName = "";
                itemRoom.ReplaceBtn = false;
                itemRoom.headIcon = "Assets/Arts/UI/BattleReady/headGray";
                itemRoom.showWeaponAction = null;
                itemRoom.showWeaponAction2 = null;
                itemRoom.showWeaponAction3 = null;
            }
        }

        var spectatorList = new List<Player>();
        for (int i = 0; i < Players.Count; i++)
        {
            if (Players[i] != null && Players[i].seat == -1)
            {
                spectatorList.Add(Players[i]);
            }
        }
        spectatorList.Sort(SortPlayer);
        for (int i = 0; i < Props.JoinList.Count; i++)
        {
            var joinItem = Props.JoinList.Get(i);

            if (spectatorList.Count > i)
            {
                joinItem.Transform?.gameObject.SetActive(true);
                joinItem.GameNameId = spectatorList[i].name;
                joinItem.Owner = spectatorList[i].Character == 0;
                joinItem.GameNameColor = spectatorList[i] == Self ? new Color(0, 1, 1) : Color.white;
            }
            else
            {
                joinItem.GameNameId = "";
                joinItem.GameNameColor = Color.white;
                joinItem.Transform?.gameObject.SetActive(false);
            }
        }

        Props.SetPropertyChange(nameof(Props.JoinList));
        Props.SetPropertyChange(nameof(Props.ReadyList));
        this.CommitProps();
        JoinG.gameObject.SetActive(Self.seat != -1);
        ReadyG.gameObject.SetActive(Self.seat != -1 && !Self.ready && !isOwner);
        UnReadyG.gameObject.SetActive(Self.seat != -1 && Self.ready && !isOwner);
        StartG.gameObject.SetActive(isOwner);
    }

    private void Regist()
    {
        TcpInstance.Regist<RoomDataS2C>(OnMsgRoomData);
        TcpInstance.Regist<PlayerUpdateS2C>(OnMsgPlayerUpdate);
        TcpInstance.Regist<GameStartS2C>(OnMsgGameStart);
        TcpInstance.Regist<GameCountStartS2C>(OnCountStart);
        TcpInstance.Regist<GameCountCancelS2C>(OnCountCancel);
        TcpInstance.Regist<RoomBreakupS2C>(OnRoomBreakup);
        TcpInstance.Regist<ExitRoomS2C>(OnMsgQuitRoom);
    }

    private void UnRegist()
    {
        TcpInstance.UnRegist<RoomDataS2C>(OnMsgRoomData);
        TcpInstance.UnRegist<PlayerUpdateS2C>(OnMsgPlayerUpdate);
        TcpInstance.UnRegist<GameStartS2C>(OnMsgGameStart);
        TcpInstance.UnRegist<GameCountStartS2C>(OnCountStart);
        TcpInstance.UnRegist<GameCountCancelS2C>(OnCountCancel);
        TcpInstance.UnRegist<RoomBreakupS2C>(OnRoomBreakup);
        TcpInstance.UnRegist<ExitRoomS2C>(OnMsgQuitRoom);
    }

    public void OnWeaponWindow(int index)
    {
        //SubWeaponWindow.SubWeaponData data = new SubWeaponWindow.SubWeaponData();
        //data.data = Self.Weapons[index];
        //data.index = index;
        //GetSubWindow<SubWeaponWindow>().Open(data);
    }

    [ButtonOnclick("Adapter/ExitButton")]
    public void QuitReady()
    {
        if (Self.Character == 0)
        {
            var mData = new Window_Tips.UIMsg_Tips();
            mData.context = "The owner will be dismissed if he exits the room";
            mData.isCancel = false;
            mData.confirm = OnQuitRoom;
            mData.title = "NOTIFY";
            UIManager.Inst.ShowWindow(WinEnum.Win_Tips, mData, true, UILayer.Top);
        }
        else
        {
            OnQuitRoom();
        }
    }

    public void OnQuitRoom()
    {
        cancellationTokenSource?.Cancel();
        isOpen = false;
        UIManager.Inst.ShowWindow(WinEnum.Win_Main);
        UIManager.Inst.CloseWindow(WinEnum.Win_BattleReady);
        TcpInstance.Send(new ExitRoomC2S(), (ExitRoomS2C msg) =>
        {
            if (msg.Code != CodePBType.Success)
            {
                Debug.Log(msg.Code);
            }
            else
            {
                UnRegist();
                TcpInstance = null;
            }
        });
    }
    [ButtonOnclick("Adapter/OBGroup/JoinButton")]
    public void OnStand()
    {
        TcpInstance.Send(new StandUpC2S(), (StandUpS2C msg) => { Debug.Log(msg.Code); });
    }

    [ButtonOnclick("Adapter/ReadyButton")]
    public void OnReady()
    {
        cancellationTokenSource = new CancellationTokenSource();
        TcpInstance.Send(new ReadyC2S { Value = true }, (ReadyS2C msg) => { Debug.Log(msg.Code); });
    }

    private bool m_Counting;

    [ButtonOnclick("Adapter/StartButton")]
    public void OnStart()
    {
        OnStartGame();
    }
    [ButtonOnclick("Adapter/CancelButton")]
    public void OnCancel()
    {

        OnStartGame();
    }
    public void OnStartGame()
    {
        cancellationTokenSource = new CancellationTokenSource();
        //bool left = false;
        //bool right = false;
        //for (int i = 0; i < Players.Count; i++)
        //{
        //    if (Players[i].seat > -1 && Players[i].seat < 3)
        //        left = true;
        //    if (Players[i].seat > 2)
        //        right = true;
        //    if (Players[i] != Self && !Players[i].ready)
        //    {
        //        Tip_Context.text = "Some players are not ready";
        //        ShowTips().Forget();
        //        return;
        //    }
        //}
        //if (!left || !right) {
        //    Tip_Context.text = "There are no players in the camp";
        //    ShowTips().Forget();
        //    return;
        //}
        TcpInstance.Send(new CountC2S { Start = !m_Counting }, (CountS2C msg) =>
        {
            Debug.Log(msg.Code);
            if (msg.Code != CodePBType.Success)
            {
                Tip_Context.text = msg.Code.ToString();
                if (msg.Code == CodePBType.NotAllReady)
                    Tip_Context.text = "Some players are not ready";
                if (msg.Code == CodePBType.NoPermission)
                    Tip_Context.text = "No Permission";
                if (msg.Code == CodePBType.PlayerNumberError)
                    Tip_Context.text = "Insufficient number of players";
                ShowTips().Forget();
            }
        });
    }
    public async UniTaskVoid ShowTips()
    {
        if (Tips.gameObject.activeInHierarchy)
            return;
        Tips.gameObject.SetActive(true);
        await UniTask.Delay(TimeSpan.FromSeconds(1));
        Tips.gameObject.SetActive(false);
    }

    [ButtonOnclick("Adapter/UNReadyButton")]
    public void OnUnReady()
    {
        TcpInstance.Send(new ReadyC2S { Value = false }, (ReadyS2C msg) => { Debug.Log(msg.Code); });
    }

    public void Arrow(int seat)
    {
        if (Self != null && Self.seat != -1)
        {
            OnMove(seat);
        }
        else
        {
            OnSit(seat);
        }
        JoinG.gameObject.SetActive(Self.seat != -1);
    }

    public void OnSit(int seat)
    {
        TcpInstance.Send(new SitDownC2S { Seat = seat }, (SitDownS2C msg) =>
        {
            Debug.Log(msg.Code);
            if (msg.Code == CodePBType.GameInCount)
            {
                Tip_Context.text = "Unable to seat at the beginning of the countdown";
                ShowTips().Forget();
            }
        });
    }

    public void OnMove(int seat)
    {
        TcpInstance.Send(new MoveC2S { Seat = seat }, (MoveS2C msg) => { Debug.Log(msg.Code); });
    }

    private void OnMsgRoomData(RoomDataS2C msg)
    {
        InitPlayer(msg.RoomData.Players);
    }

    private void OnMsgPlayerUpdate(PlayerUpdateS2C msg)
    {
        switch (msg.Mode)
        {
            case UpdatePBType.Add:
                Players.Add(new Player { name = msg.Player.Name, seat = msg.Player.Seat, Character = (ulong)Players.Count,Weapons=msg.Player.Weapons });
                break;
            case UpdatePBType.Remove:
                Players.RemoveAll(m => m.name == msg.Player.Name);
                break;
            case UpdatePBType.Update:
                var target = Players.Find(m => m.name == msg.Player.Name);
                target.seat = msg.Player.Seat;
                target.ready = msg.Player.Ready;
                target.Weapons=msg.Player.Weapons;
                break;
        }

        UpdateSeats();
    }

    private void OnMsgQuitRoom(ExitRoomS2C msg)
    {
        if (Self.Character != 0)
        {
            bool notify = true;
            for (int i = 0; i < Players.Count; i++)
            {
                if (Players[i].Character == 0)
                    notify = false;
            }
            if (!notify)
                return;

        }
    }

    private void OnCountStart(GameCountStartS2C msg)
    {
        m_Counting = true;
        WaitTime = msg.Scend;
        ReadyAll().Forget();
        if (isOwner)
        {
            StartG.gameObject.SetActive(false);
            CancelG.gameObject.SetActive(true);
        }
    }

    private void OnCountCancel(GameCountCancelS2C _)
    {
        m_Counting = false;
        cancellationTokenSource?.Cancel();
        CutDownImg.gameObject.SetActive(false);
        if (isOwner)
        {
            StartG.gameObject.SetActive(true);
            CancelG.gameObject.SetActive(false);
        }

    }

    private void OnRoomBreakup(RoomBreakupS2C _)
    {
        // 房间解散
        var mData = new Window_Tips.UIMsg_Tips();
        mData.context = "Room has been disbanded";
        mData.isCancel = false;
        mData.confirm = OnQuitRoom;
        mData.title = "NOTIFY";
        UIManager.Inst.ShowWindow(WinEnum.Win_Tips, mData, true, UILayer.Top);
    }

    private void OnMsgGameStart(GameStartS2C msg)
    {
        m_Counting = false;
        TestBattle.StartMsg = msg;

        var dic = new Dictionary<int, string>();
        foreach (var player in Players)
            dic.Add(player.seat, player.name);
        TestBattle.PlayerNameDic = dic;
        var dic2 = new Dictionary<int, ulong[]>();
        foreach (var player in Players)
            dic2.Add(player.seat, player.Weapons);
        TestBattle.PlayerWeaponDic = dic2;
        PlayerPrefs.SetInt("Test", 0);
        cancellationTokenSource?.Cancel();
        CutDownImg.gameObject.SetActive(false);
        ScreenViewManager.Inst.MainLayer.BeginNavTo(ScreenViewEnum.Battle);
        UIManager.Inst.ScalMinWindow(WinEnum.Win_BattleReady);
    }

    public override void Close()
    {
        isOpen = false;
        base.Close();
    }

    public override void Destroy()
    {
        base.Destroy();
        if (TcpInstance != null)
            UnRegist();
    }

    public int SortPlayer(Player A, Player B)
    {
        if (A.Character == 0) return 0;
        if (A == Self) return 1;
        return A.Character.CompareTo(B.Character);
    }
}