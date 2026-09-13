// 休憩タスクの名前。シフト表・タスク一覧の表記、およびAPI側のbreakTaskName
// (api/lib/usecase/shift_usecase.go)と一致させる必要がある
const String breakTaskName = '休憩';

// シフトカードのデータモデル
class ShiftMember {
  final String name;
  final String grade;
  final String bureau;

  ShiftMember({
    required this.name,
    required this.grade,
    required this.bureau,
  });
}
class ShiftMembers {
  final String sTime;
  final String eTime;
  final List<ShiftMember> members;

  ShiftMembers({
    required this.sTime,
    required this.eTime,
    required this.members,
  });
}

class ShiftCardData {
  final String taskName;
  final String startTime;
  final String endTime;
  final String place;
  final String url;
  final String manualUrl;
  final List<ShiftMembers> shiftMembers;
  final ShiftMembers beforeMembers;
  final ShiftMembers afterMembers;

  // 休憩は通常のシフトと扱いを変える。集合場所・マニュアル・担当者一覧を持たず、
  // Newバッジもレビューも出さない(#488)
  bool get isBreak => taskName.trim() == breakTaskName;

  ShiftCardData({
    required this.taskName,
    required this.startTime,
    required this.endTime,
    required this.place,
    required this.url,
    required this.manualUrl,
    required this.shiftMembers,
    required this.beforeMembers,
    required this.afterMembers,
  });
}


class ShiftCardDataList {
  final List<ShiftCardData> data;
  ShiftCardDataList(this.data);
  
  // jsonからShiftCardDataListを生成するファクトリコンストラクタ
  factory ShiftCardDataList.fromJson(List<dynamic> json) {
    List<ShiftCardData> data = json.map((item) => ShiftCardData(
      taskName: item['task_name'] as String,
      startTime: item['start_time'] as String,
      endTime: item['end_time'] as String,
      place: item['place'] as String,
      url: item['url'] as String,
      manualUrl: (item['manual_url'] as String?) ?? '',
      shiftMembers: item['shift_members'] != null?
        (item['shift_members'] as List<dynamic>)
          .map((member) => ShiftMembers(
                sTime: member['s_time'],
                eTime: member['e_time'],
                members: member['members'] != null?
                  (member['members'] as List<dynamic>)
                    .map((m) => ShiftMember(
                          name: m['name'],
                          grade: m['grade'],
                          bureau: m['bureau'],
                        ))
                    .toList():
                  [],
              ))
          .toList():
        [],
      beforeMembers: ShiftMembers(
        sTime: item['before_members']['s_time'],
        eTime: item['before_members']['e_time'],
        members: item['before_members']['members'] != null?
          (item['before_members']['members'] as List<dynamic>)
            .map((m) => ShiftMember(
                name: m['name']?? 'データの取得に失敗しました',
                grade: m['grade']?? '',
                bureau: m['bureau']?? '',
            ))
            .toList(): 
          [],
      ),
      afterMembers: ShiftMembers(
        sTime: item['after_members']['s_time'],
        eTime: item['after_members']['e_time'],
        members: item['after_members']['members'] != null?
          (item['after_members']['members'] as List<dynamic>)
            .map((m) => ShiftMember(
                name: m['name']?? 'データの取得に失敗しました',
                grade: m['grade']?? '',
                bureau: m['bureau']?? '',
            ))
            .toList():
          [],
      ),
    )).toList();
    return ShiftCardDataList(data);
  }
}