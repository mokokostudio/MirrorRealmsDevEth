using MR.Battle;
using TMPro;
using UnityEngine;
using UnityEngine.UI;

public class HPBar : MonoBehaviour {
    public UnitCD Unit { get; set; }

    public TMP_Text nameTxt1;
    public TMP_Text nameTxt2;
    public Slider HPSlider1;
    public Slider HPSlider2;
    public Slider SPSlider;

    private bool m_ShowSP;

    private bool m_CampFlag;

    private Slider m_HPSlider;

    private void Start() {
        nameTxt1.text = nameTxt2.text = TestBattle.PlayerNameDic[Unit.Player.Index];
        m_ShowSP = true;
        UpdateCamp();
    }

    private void Update() {
        transform.rotation = Camera.main.transform.rotation;
        if (m_ShowSP)
            SPSlider.value = Unit.Player.SP.AsFloat() / Config.Battle.Constant.SPMax;
        if (m_CampFlag != (Battle.Instance.CameraPlayer.Unit.Camp == Unit.Camp))
            UpdateCamp();
        m_HPSlider.value = Unit.HP.AsFloat() / Config.Battle.Constant.HPMax;
    }

    private void UpdateCamp() {
        if (Battle.Instance.CameraPlayer.Unit.Camp == Unit.Camp) {
            nameTxt1.gameObject.SetActive(true);
            nameTxt2.gameObject.SetActive(false);
            m_HPSlider = HPSlider1;
            HPSlider1.gameObject.SetActive(true);
            HPSlider2.gameObject.SetActive(false);
        } else {
            nameTxt1.gameObject.SetActive(false);
            nameTxt2.gameObject.SetActive(true);
            m_HPSlider = HPSlider2;
            HPSlider1.gameObject.SetActive(false);
            HPSlider2.gameObject.SetActive(true);
        }
        m_CampFlag = Battle.Instance.CameraPlayer.Unit.Camp == Unit.Camp;
    }
}
