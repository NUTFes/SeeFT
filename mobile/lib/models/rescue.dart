// レスキュー画面のタスクのドロップダウンの選択肢として使用するデータモデル
class RescueTaskDropdownMenuItem {
  final int id;
  final String taskName;

  RescueTaskDropdownMenuItem({
    required this.id,
    required this.taskName
  });
}

// ===== レスキューの返答を表すデータモデル ======
// 共通の基底クラス
abstract class RescueResponse {
  final String type;
  final int id;
  final String userName;
  final String time;
  final String status;
  final String response;

  RescueResponse({
    required this.type,
    required this.id,
    required this.userName,
    required this.time,
    required this.status,
    required this.response,
  });

  // factory RescueResponse.fromJson(Map<String, dynamic> json) {
  factory RescueResponse.fromJson(dynamic json) {
    switch (json['type']) {
      case 'trouble':
        return TroubleRescueResponse.fromJson(json);
      case 'question':
        return QuestionRescueResponse.fromJson(json);
      case 'shorthanded':
        return ShorthandedRescueResponse.fromJson(json);
      default:
        throw Exception('Unknown type: ${json['type']}');
    }
  }

  Map<String, dynamic> toJson();
}

// --- Trouble ---
class TroubleContent {
  final String task;
  final String place;
  final String detail;

  TroubleContent({
    required this.task,
    required this.place,
    required this.detail,
  });
  
  // factory TroubleContent.fromJson(Map<String, dynamic> json) {
  factory TroubleContent.fromJson(dynamic json) {
    return TroubleContent(
      task: json['task'],
      place: json['place'],
      detail: json['detail'],
    );
  }

  Map<String, dynamic> toJson() => {
        'task': task,
        'place': place,
        'detail': detail,
      };
}

class TroubleRescueResponse extends RescueResponse {
  final TroubleContent content;

  TroubleRescueResponse({
    required int id,
    required String userName,
    required String time,
    required String status,
    required String response,
    required this.content,
  }) : super(
          type: 'trouble',
          id: id,
          userName: userName,
          time: time,
          status: status,
          response: response,
        );

  // factory TroubleRescueResponse.fromJson(Map<String, dynamic> json) {
  factory TroubleRescueResponse.fromJson(dynamic json) {
    print(json);
    return TroubleRescueResponse(
      id: json['id'],
      userName: json['user_name'],
      time: json['time'],
      status: json['status'],
      response: json['response'],
      content: TroubleContent.fromJson(json['content']),
    );
  }

  @override
  Map<String, dynamic> toJson() => {
        'type': type,
        'id': id,
        'user_name': userName,
        'time': time,
        'status': status,
        'response': response,
        'content': content.toJson(),
      };
}

// --- Question ---
class QuestionContent {
  final String question;

  QuestionContent({required this.question});

  // factory QuestionContent.fromJson(Map<String, dynamic> json) {
  factory QuestionContent.fromJson(dynamic json) {
    return QuestionContent(
      question: json['question'],
    );
  }

  Map<String, dynamic> toJson() => {
        'question': question,
      };
}

class QuestionRescueResponse extends RescueResponse {
  final QuestionContent content;

  QuestionRescueResponse({
    required int id,
    required String userName,
    required String time,
    required String status,
    required String response,
    required this.content,
  }) : super(
          type: 'question',
          id: id,
          userName: userName,
          time: time,
          status: status,
          response: response,
        );

  // factory QuestionRescueResponse.fromJson(Map<String, dynamic> json) {
  factory QuestionRescueResponse.fromJson(dynamic json) {
    return QuestionRescueResponse(
      id: json['id'],
      userName: json['user_name'],
      time: json['time'],
      status: json['status'],
      response: json['response'],
      content: QuestionContent.fromJson(json['content']),
    );
  }

  @override
  Map<String, dynamic> toJson() => {
        'type': type,
        'id': id,
        'user_name': userName,
        'time': time,
        'status': status,
        'response': response,
        'content': content.toJson(),
      };
}

// --- Shorthanded ---
class ShorthandedContent {
  final String task;
  final int missingNumber;
  final String place;

  ShorthandedContent({
    required this.task,
    required this.missingNumber,
    required this.place,
  });

  // factory ShorthandedContent.fromJson(Map<String, dynamic> json) {
  factory ShorthandedContent.fromJson(dynamic json) {
    return ShorthandedContent(
      task: json['task'],
      missingNumber: json['missing_number'],
      place: json['place'],
    );
  }

  Map<String, dynamic> toJson() => {
        'task': task,
        'missing_number': missingNumber,
        'place': place,
      };
}

class ShorthandedRescueResponse extends RescueResponse {
  final ShorthandedContent content;

  ShorthandedRescueResponse({
    required int id,
    required String userName,
    required String time,
    required String status,
    required String response,
    required this.content,
  }) : super(
          type: 'shorthanded',
          id: id,
          userName: userName,
          time: time,
          status: status,
          response: response,
        );

  // factory ShorthandedRescueResponse.fromJson(Map<String, dynamic> json) {
  factory ShorthandedRescueResponse.fromJson(dynamic json) {
    return ShorthandedRescueResponse(
      id: json['id'],
      userName: json['user_name'],
      // time: DateUtil.parseCustomDate(json['time']), // 共通パーサーを使用
      time: json['time'], // 共通パーサーを使用
      status: json['status'],
      response: json['response'],
      content: ShorthandedContent.fromJson(json['content']),
    );
  }

  @override
  Map<String, dynamic> toJson() => {
        'type': type,
        'id': id,
        'user_name': userName,
        'time': time,
        'status': status,
        'response': response,
        'content': content.toJson(),
      };
}


// // 日付のパースを行うユーティリティクラス
// class DateUtil {
//   static DateTime parseCustomDate(String? dateString) {
//     if (dateString == null || dateString.isEmpty) {
//       return DateTime.now();
//     }
    
//     try {
//       // 複数の形式に対応
//       if (dateString.contains('/')) {
//         // "2025/08/09 19:31:48" → "2025-08-09 19:31:48"
//         String normalizedDate = dateString.replaceAll('/', '-');
//         return DateTime.parse(normalizedDate);
//       } else {
//         // 通常のISO形式
//         return DateTime.parse(dateString);
//       }
//     } catch (e) {
//       print('Date parsing error: $e for date string: $dateString');
//       return DateTime.now();
//     }
//   }
// }