import 'package:flutter/material.dart';
import 'package:seeft_mobile/models/shift_card.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';
import 'package:seeft_mobile/theme/tokens.dart';

class ShiftCardManualDemoPage extends StatelessWidget {
  const ShiftCardManualDemoPage({super.key});

  @override
  Widget build(BuildContext context) {
    // 各マニュアルのシフトカードデモデータ
    final demoCards = <ShiftCardData>[
      _card('電力配線', '08:00', '12:00', 'B講義室', '/manuals/haisen.html'),
      _card('駐車場設営・撤収', '08:00', '10:00', '第2駐車場', '/manuals/parking.html'),
      _card('案内所準備・片付け', '09:00', '10:00', '電気1号棟1階103', '/manuals/annai.html'),
      _card('本部設営', '09:00', '11:00', '電気1号棟1階 104', '/manuals/honbu.html'),
      _card('のぼり広告片付け', '16:00', '17:00', '講義棟前', '/manuals/nobori.html'),
      _card('物販テント運営', '10:00', '17:00', '物販テントエリア', '/manuals/buppan.html'),
      _card('幼稚園WARSコラボブース', '10:00', '17:00', 'AL3', '/manuals/wars.html'),
      _card('お化け屋敷', '10:00', '16:00', '物材棟2F 大学院講義室', '/manuals/obake.html'),
    ];

    return Scaffold(
      backgroundColor: AppColors.base,
      appBar: AppBar(
        title: const Text('マイシフト'),
        backgroundColor: AppColors.main,
        foregroundColor: AppColors.textWhite,
      ),
      body: ListView.separated(
        padding: const EdgeInsets.all(16.0),
        itemCount: demoCards.length,
        separatorBuilder: (_, __) => const SizedBox(height: 12),
        itemBuilder: (context, index) => ShiftCard(data: demoCards[index]),
      ),
    );
  }

  static ShiftCardData _card(
    String task, String start, String end, String place, String url,
  ) {
    return ShiftCardData(
      taskName: task,
      startTime: start,
      endTime: end,
      place: place,
      url: url,
      shiftMembers: [
        ShiftMembers(
          s_time: start,
          e_time: end,
          members: [
            ShiftMember(name: '技大太郎', grade: '3', bureau: '総務局'),
          ],
        ),
      ],
      beforeMembers: ShiftMembers(s_time: '', e_time: '', members: []),
      afterMembers: ShiftMembers(s_time: '', e_time: '', members: []),
    );
  }
}
