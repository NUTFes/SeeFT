import 'package:flutter/material.dart';
import 'package:seeft_mobile/theme/tokens.dart';
import 'package:seeft_mobile/configs/importer.dart';
import 'package:url_launcher/url_launcher.dart';

// // 画像のUIを再現したウィジェット
// class TaskCard extends StatelessWidget {
//   const TaskCard({super.key});

//   @override
//   Widget build(BuildContext context) {
//     return Card(
//       margin: const EdgeInsets.all(0),
//       shape: RoundedRectangleBorder(
//         borderRadius: BorderRadius.circular(12.0),
//       ),
//       elevation: 4.0,
//       child: Column(
//         mainAxisSize: MainAxisSize.min, // Columnが必要な分だけ高さを取るように
//         crossAxisAlignment: CrossAxisAlignment.start,
//         children: [
//           // 1. 上部のヘッダー（時間とボタン）
//           _buildHeader(),
          
//           Padding(
//             padding: const EdgeInsets.fromLTRB(8.0, 6.0, 8.0, 8.0),
//             child: Column(
//               crossAxisAlignment: CrossAxisAlignment.start,
//               children: [
//                 // 2. メインタイトル
//                 const Text(
//                   '参加団体窓口対応',
//                   style: TextStyle(fontSize: AppFontSizes.md, fontWeight: FontWeight.bold, color: AppColors.textBlack),
//                 ),
//                 const SizedBox(height: 6.0),
//                 const Divider(
//                   height: 1,
//                   color: AppColors.grayLight, // 区切り線の色
//                 ),
//                 const SizedBox(height: 6.0),
//                 // 3. 区切り線と詳細情報（開閉式）
//                 Theme(
//                   data: Theme.of(context).copyWith(dividerColor: Colors.transparent), // 区切り線の色を透明に
//                   child: _buildDetailsSection(),
//                 ),
//               ],
//             ),
//           ),
//         ],
//       ),
//     );
//   }

//   // ヘッダー部分を生成するヘルパーメソッド
//   Widget _buildHeader() {
//     return Padding(
//       padding: const EdgeInsets.fromLTRB(8.0, 8.0, 8.0, 8.0),
//       child: Row(
//         mainAxisAlignment: MainAxisAlignment.spaceBetween,
//         children: [
//           const Row(
//             children: [
//               Icon(Icons.access_time, color: AppColors.textBlack, size: 16),
//               SizedBox(width: 2.0),
//               Text(
//                 '8:00〜9:00',
//                 style: TextStyle(fontSize: AppFontSizes.sm, fontWeight: FontWeight.bold, color: AppColors.textBlack),
//               ),
//             ],
//           ),
//           ElevatedButton(
//             onPressed: () {
//               // マニュアルを開くボタンの処理
//             },
//             style: ElevatedButton.styleFrom(
//               backgroundColor: AppColors.sub,
//               foregroundColor: AppColors.textWhite,
//               shape: const StadiumBorder(), // 角が丸い形状
//               elevation: 0, // 影を消す
//               padding: const EdgeInsets.symmetric(horizontal: 8.0, vertical: 4.0),
//             ),
//             child: const Text(
//               'マニュアルを開く',
//               style: TextStyle(fontSize: AppFontSizes.xs, color: AppColors.textWhite),
//             ),
//           ),
//         ],
//       ),
//     );
//   }

//   // 詳細情報部分（開閉式）を生成するヘルパーメソッド
//   Widget _buildDetailsSection() {
//     // ExpansionTileを使うと、簡単に開閉できるUIが作れます
//     return ExpansionTile(
//       // タイルの余白を消してデザインを調整
//       tilePadding: EdgeInsets.zero,
//       childrenPadding: EdgeInsets.zero,
//       // デフォルトで表示されるタイトル部分
//       title: const Row(
//         children: [
//           Icon(Icons.location_on_outlined, color: AppColors.textBlack, size: 16,),
//           SizedBox(width: 2.0),
//           Text(
//             '物材棟調理場(1〜5F)',
//             style: TextStyle(fontSize: AppFontSizes.sm, color: AppColors.textBlack),
//           ),
//         ],
//       ),
//       // trailing: const Icon(Icons.expand_more, color: AppColors.textBlack), // アイコンの色を変更
      
//       // タイルが開いた時の背景色
//       // backgroundColor: AppColors.grayLight,
//       // タイルが開いた時のテキストスタイル
//       // textColor: AppColors.textBlack,
//       // アイコンの色を変更
//       // iconColor: AppColors.textBlack,
      
//       // 開いた時に表示されるコンテンツ
//       children: <Widget>[
//         // const Divider(height: 1),
//         const SizedBox(height: 16),
//         _buildSection(
//           title: '【集合場所】',
//           content: const ['物材棟調理場(1〜5F)'],
//         ),
//         _buildSection(
//           title: '【困った時は】',
//           content: [
//             '1. 近くの先輩に聞く',
//             '2. それでも分からなかったら本部に連絡',
//           ],
//         ),
//         _buildSection(
//           title: '【担当者の一覧】',
//           content: [
//             '8:00〜8:15',
//             '(情報局B3)技大一郎, (情報局B3)技大次郎',
//             '\n8:15〜8:30', // 少し間を空けるために改行
//             '(情報局B3)技大次郎',
//           ],
//         ),
//         _buildSection(
//           title: '【前の時間の担当者の一覧】',
//           content: [
//             '7:45〜8:00',
//             '(情報局B3)技大一郎',
//           ],
//         ),
//         _buildSection(
//           title: '【次の時間の担当者の一覧】',
//           content: [
//             '8:30〜8:45',
//             '(情報局B3)技大三郎',
//           ],
//         ),
//       ],
//     );
//   }
  
//   // 各セクションのUIを生成するヘルパーメソッド
//   Widget _buildSection({required String title, required List<String> content}) {
//     return Padding(
//       padding: const EdgeInsets.only(bottom: 16.0),
//       child: Align(
//         alignment: Alignment.centerLeft,
//         child: Column(
//           crossAxisAlignment: CrossAxisAlignment.start,
//           children: [
//             Text(
//               title,
//               style: const TextStyle(fontWeight: FontWeight.bold),
//             ),
//             const SizedBox(height: 4.0),
//             ...content.map((line) => Text(line, style: const TextStyle(height: 1.5))),
//           ],
//         ),
//       ),
//     );
//   }
// }

// // 担当者情報をまとめるクラス
// class StaffInfo {
//   final String time;
//   final List<String> names;

//   StaffInfo({required this.time, required this.names});
// }

// // カード全体のデータをまとめるモデルクラス
// class TaskCardData {
//   final String timeRange;
//   final String title;
//   final String location;
//   final String meetingPoint;
//   final List<String> helpSteps;
//   final List<StaffInfo> currentStaff;
//   final StaffInfo? previousStaff; // 前の担当者はいない場合もあるので `?` をつけてnull許容に
//   final StaffInfo? nextStaff;     // 次の担当者も同様

//   TaskCardData({
//     required this.timeRange,
//     required this.title,
//     required this.location,
//     required this.meetingPoint,
//     required this.helpSteps,
//     required this.currentStaff,
//     this.previousStaff,
//     this.nextStaff,
//   });
// }


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
  final String s_time;
  final String e_time;
  final List<ShiftMember> members;

  ShiftMembers({
    required this.s_time,
    required this.e_time,
    required this.members,
  });
}

class ShiftCardData {
  final String taskName;
  final String startTime;
  final String endTime;
  final String place;
  final String url;
  final List<ShiftMembers> shiftMembers;
  final ShiftMembers beforeMembers;
  final ShiftMembers afterMembers;

  ShiftCardData({
    required this.taskName,
    required this.startTime,
    required this.endTime,
    required this.place,
    required this.url,
    required this.shiftMembers,
    required this.beforeMembers,
    required this.afterMembers,
  });
}
// 画像のUIを再現したウィジェット
class ShiftCard extends StatelessWidget {
  // finalでデータモデルを受け取る
  final ShiftCardData data;

  const ShiftCard({super.key, required this.data});
  
  // リンクを開くための非同期メソッドを定義
  Future<void> _launchManualUrl(String url) async {
    // 開きたいURLをUriオブジェクトに変換
    final Uri uri = Uri.parse(url);

    // 3. launchUrlを実行
    if (!await launchUrl(uri, mode: LaunchMode.externalApplication)) {  // ここで外部アプリで開くように指定
      // 4. URLが開けなかった場合のエラー処理
      throw Exception('Could not launch $uri');
    }
  }

  @override
  Widget build(BuildContext context) {
    // UIの構造は同じだが、表示するテキストは `data` オブジェクトから取得する
    return Card(
      margin:EdgeInsets.zero,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8.0),
        side: BorderSide(
          color: AppColors.grayLight, // 枠線の色
          width: 1.0, // 枠線の太さ
        ),
      ),
      elevation: 0.5,
      shadowColor: null,
      color: AppColors.base, // 背景色を設定
      // shape: RoundedRectangleBorder(
      //   borderRadius: BorderRadius.circular(8.0),
      // ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // Header
          Padding(
            padding: const EdgeInsets.fromLTRB(8.0, 8.0, 8.0, 3.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Row(
                  children: [
                    const Icon(
                      Icons.access_time,
                      color: AppColors.textBlack,
                      size: 16
                    ),
                    const SizedBox(width: 2.0),
                    // データモデルから時間を表示
                    Text(
                      data.startTime + "〜" + data.endTime,
                      style: const TextStyle(
                        fontSize: AppFontSizes.sm,
                        color: AppColors.textBlack
                      ),
                    ),
                  ],
                ),
                ElevatedButton(
                  onPressed: () {
                    // マニュアルを開くボタンの処理
                    if (data.url.isNotEmpty) {
                      _launchManualUrl(data.url);
                    } else {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('マニュアルのURLが設定されていません')),
                      );
                    }
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.sub,
                    foregroundColor: AppColors.textWhite,
                    shape: const StadiumBorder(), // 角が丸い形状
                    elevation: 0, // 影を消す
                    padding: const EdgeInsets.symmetric(horizontal: 8.0, vertical: 4.0),
                    // minimumSize: const Size(0, 22), // 最小サイズを0に設定
                    // fixedSize: const Size(112, 22), // 固定サイズを0に設定
                  ),
                  child: const Text(
                    'マニュアルを開く',
                    style: TextStyle(
                      fontSize: AppFontSizes.xs, 
                      color: AppColors.textWhite
                    ),
                  ),
                ),
              ],
            ),
          ),
          
          Padding(
            padding: const EdgeInsets.fromLTRB(8.0, 3.0, 8.0, 8.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // データモデルからタイトルを表示
                Text(
                  data.taskName,
                  style: const TextStyle(
                    fontSize: AppFontSizes.md,
                    color: AppColors.textBlack, 
                    fontWeight: FontWeight.bold
                  ),
                ),
                const SizedBox(height: 6.0),
                const Divider(
                  height: 1,
                  color: AppColors.grayLight, // 区切り線の色
                ),
                // const SizedBox(height: 6.0),
                Theme(
                  data: Theme.of(context).copyWith(dividerColor: Colors.transparent), // 区切り線の色を透明に
                  child: ExpansionTile(
                    tilePadding: EdgeInsets.zero,
                    iconColor: AppColors.textBlack, // アイコンの色を変更
                    minTileHeight: 0, // タイルの高さを最小に
                    // データモデルから場所を表示
                    title: Row(
                      children: [
                        const Icon(
                          Icons.location_on_outlined, 
                          color: AppColors.textBlack,
                          size: 16
                        ),
                        const SizedBox(width: 2.0),
                        Text(
                          data.place,
                          style: const TextStyle(
                            fontSize: AppFontSizes.sm,
                            color: AppColors.textBlack
                          ),
                        ),
                      ],
                    ),
                    children: <Widget>[
                      // const SizedBox(height: 6.0),
                      _buildSection(
                        title: '【集合場所】',
                        content: [data.place]
                      ),
                      _buildSection(
                        title: '【困った時は】',
                        content: [
                          '1. 近くの先輩に聞く',
                          '2. それでも分からなかったら本部に連絡',
                        ],
                      ),
                      _buildSection(
                        title: '【担当者の一覧】',
                        content: [
                          for (var member in data.shiftMembers)
                            '${member.s_time}〜${member.e_time}\n' +
                            member.members.map((m) => '(${m.bureau}${m.grade}) ${m.name}').join(', '),
                        ],
                      ),
                      _buildSection(
                        title: '【前の時間の担当者の一覧】',
                        content: [
                          if (data.beforeMembers.members.isNotEmpty)
                            '${data.beforeMembers.s_time}〜${data.beforeMembers.e_time}\n' +
                            data.beforeMembers.members.map((m) => '(${m.bureau}${m.grade}) ${m.name}').join(', ')
                          else
                            '前の時間の担当者はいません',
                        ],
                      ),
                      _buildSection(
                        title: '【次の時間の担当者の一覧】',
                        content: [
                          if (data.afterMembers.members.isNotEmpty)
                            '${data.afterMembers.s_time}〜${data.afterMembers.e_time}\n' +
                            data.afterMembers.members.map((m) => '(${m.bureau}${m.grade}) ${m.name}').join(', ')
                          else
                            '次の時間の担当者はいません',
                        ],
                      ),
                    ],
                  ),
                )
              ],
            ),
          ),
        ],
      ),
    );
  }
  // 各セクションのUIを生成するヘルパーメソッド
  Widget _buildSection({required String title, required List<String> content}) {
    return Padding(
      padding: const EdgeInsets.only(top: 8.0, bottom: 8.0),
      child: Align(
        alignment: Alignment.centerLeft,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: const TextStyle(
                fontSize: AppFontSizes.xs,
                color: AppColors.textBlack,
                fontWeight: FontWeight.bold
              ),
            ),
            const SizedBox(height: 4.0),
            ...content.map((line) => Text(
              line, 
              style: const TextStyle(
                fontSize: AppFontSizes.xs, 
                color: AppColors.textBlack,
                height: 1.5 // 行間を調整
              )
            )),
          ],
        ),
      ),
    );
  }
}

// // 画像のUIを再現したウィジェット
// class TaskCardRefactored extends StatelessWidget {
//   // finalでデータモデルを受け取る
//   final TaskCardData data;

//   const TaskCardRefactored({super.key, required this.data});

//   @override
//   Widget build(BuildContext context) {
//     // UIの構造は同じだが、表示するテキストは `data` オブジェクトから取得する
//     return Card(
//       // margin: const EdgeInsets.all(16.0),
//       margin: const EdgeInsets.all(0.0),
//       shape: RoundedRectangleBorder(
//         borderRadius: BorderRadius.circular(8.0),
//       ),
//       child: Column(
//         mainAxisSize: MainAxisSize.min,
//         children: [
//           // Header
//           Padding(
//             padding: const EdgeInsets.fromLTRB(16.0, 8.0, 8.0, 8.0),
//             child: Row(
//               mainAxisAlignment: MainAxisAlignment.spaceBetween,
//               children: [
//                 Row(
//                   children: [
//                     const Icon(Icons.access_time, color: Colors.black54, size: 20),
//                     const SizedBox(width: 8.0),
//                     // データモデルから時間を表示
//                     Text(
//                       data.timeRange,
//                       style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w500),
//                     ),
//                   ],
//                 ),
//                 // ... ボタン部分は省略 ...
//               ],
//             ),
//           ),
          
//           Padding(
//             padding: const EdgeInsets.fromLTRB(16.0, 0, 16.0, 16.0),
//             child: Column(
//               crossAxisAlignment: CrossAxisAlignment.start,
//               children: [
//                 // データモデルからタイトルを表示
//                 Text(
//                   data.title,
//                   style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
//                 ),
//                 const SizedBox(height: 8.0),
                
//                 // ExpansionTile
//                 ExpansionTile(
//                   tilePadding: EdgeInsets.zero,
//                   // データモデルから場所を表示
//                   title: Row(
//                     children: [
//                       const Icon(Icons.location_on_outlined, color: Colors.black54),
//                       const SizedBox(width: 8.0),
//                       Text(
//                         data.location,
//                         style: const TextStyle(fontSize: 16),
//                       ),
//                     ],
//                   ),
//                   children: <Widget>[
//                     // ... 内部のコンテンツも同様に data オブジェクトから生成する
//                   ],
//                 ),
//               ],
//             ),
//           ),
//         ],
//       ),
//     );
//   }
// }