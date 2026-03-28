import 'package:flutter/material.dart';
import 'package:seeft_mobile/models/shift_card.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';
import 'package:seeft_mobile/theme/tokens.dart';

class ShiftCardManualDemoPage extends StatelessWidget {
  const ShiftCardManualDemoPage({super.key});

  @override
  Widget build(BuildContext context) {
    const pdfUrl = String.fromEnvironment(
      'MANUAL_PDF_URL',
      defaultValue: '',
    );

    final demoData = ShiftCardData(
      taskName: '電力配線',
      startTime: '18:00',
      endTime: '22:00',
      place: 'B講義室',
      url: pdfUrl,
      shiftMembers: [
        ShiftMembers(
          s_time: '18:00',
          e_time: '22:00',
          members: [
            ShiftMember(name: '井上英明', grade: '3', bureau: '総務局'),
            ShiftMember(name: '坪内創', grade: '3', bureau: '総務局'),
            ShiftMember(name: '小日向風磨', grade: '2', bureau: '総務局'),
            ShiftMember(name: '沓掛正太郎', grade: '3', bureau: '総務局'),
          ],
        ),
      ],
      beforeMembers: ShiftMembers(
        s_time: '',
        e_time: '',
        members: [],
      ),
      afterMembers: ShiftMembers(
        s_time: '',
        e_time: '',
        members: [],
      ),
    );

    return Scaffold(
      backgroundColor: AppColors.base,
      appBar: AppBar(
        title: const Text('マイシフト'),
        backgroundColor: AppColors.main,
        foregroundColor: AppColors.textWhite,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: ShiftCard(data: demoData),
      ),
    );
  }
}
