import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/widgets/manual_viewer.dart';
import 'package:seeft_mobile/widgets/new_badge.dart';
import 'package:seeft_mobile/widgets/review_bottom_sheet.dart';
import 'package:url_launcher/url_launcher.dart';

// シフトカードのウィジェット
class ShiftCard extends StatelessWidget {
  final ShiftCardData data;
  final int userID;
  final bool isNew;
  final VoidCallback? onOpened;

  const ShiftCard({
    super.key,
    required this.data,
    required this.userID,
    this.isNew = false,
    this.onOpened,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: EdgeInsets.zero,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(AppBorderRadius.normal),
        side: const BorderSide(
          color: AppColors.grayLight,
          width: 1.0,
        ),
      ),
      elevation: 0.5,
      shadowColor: null,
      color: AppColors.base,
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Row(
                  children: [
                    const Icon(
                      Icons.access_time,
                      color: AppColors.textBlack,
                      size: 16,
                    ),
                    const SizedBox(width: 2.0),
                    Text(
                      '${data.startTime}~${data.endTime}',
                      style: const TextStyle(
                        fontSize: AppFontSizes.sm,
                        color: AppColors.textBlack,
                      ),
                    ),
                  ],
                ),
                _cardMenu(context),
              ],
            ),
            Row(
              children: [
                if (isNew) ...[
                  const NewBadge(),
                  const SizedBox(width: 4.0),
                ],
                Flexible(
                  child: Text(
                    data.taskName.toString(),
                    style: const TextStyle(
                      fontSize: AppFontSizes.md,
                      color: AppColors.textBlack,
                      fontWeight: FontWeight.bold,
                    ),
                    textAlign: TextAlign.left,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 6.0),
            const Divider(
              height: 1,
              color: AppColors.grayLight,
            ),
            const SizedBox(height: 6.0),
            Row(
              children: [
                const Icon(
                  Icons.location_on_outlined,
                  color: AppColors.textBlack,
                  size: 16,
                ),
                const SizedBox(width: 2.0),
                Flexible(
                  child: Text(
                    data.place,
                    style: const TextStyle(
                      fontSize: AppFontSizes.sm,
                      color: AppColors.textBlack,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 4.0),
            _ManualToggle(
              url: data.url,
              onOpened: () {
                if (isNew && onOpened != null) {
                  onOpened!();
                }
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _cardMenu(BuildContext context) {
    return Builder(
      builder: (menuContext) {
        return _ShiftCardMoreButton(
          onPressed: (position) => _showCardMenu(menuContext, position),
        );
      },
    );
  }

  Future<void> _showCardMenu(BuildContext context, Offset position) async {
    final overlay = Overlay.of(context).context.findRenderObject() as RenderBox;
    final action = await showMenu<_ShiftCardMenuAction>(
      context: context,
      color: Colors.transparent,
      elevation: 0,
      shadowColor: Colors.transparent,
      surfaceTintColor: Colors.transparent,
      constraints: const BoxConstraints.tightFor(width: 212),
      menuPadding: EdgeInsets.zero,
      position: RelativeRect.fromRect(
        Rect.fromLTWH(position.dx, position.dy, 0, 0),
        Offset.zero & overlay.size,
      ),
      items: const [
        _ShiftCardMenuPanel(),
      ],
    );

    if (action == null) return;
    switch (action) {
      case _ShiftCardMenuAction.members:
        if (isNew && onOpened != null) {
          onOpened!();
        }
        if (context.mounted) {
          _showMembersDialog(context);
        }
        break;
      case _ShiftCardMenuAction.review:
        if (context.mounted) {
          ReviewBottomSheet.show(context, data.taskName, userID);
        }
        break;
    }
  }

  void _showMembersDialog(BuildContext context) {
    showDialog(
      context: context,
      barrierColor: Colors.black.withValues(alpha: 0.2),
      builder: (context) {
        return AlertDialog(
          backgroundColor: AppColors.base,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppBorderRadius.normal),
          ),
          contentPadding: const EdgeInsets.all(16.0),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildMemberSection(
                  title: '【担当者の一覧】',
                  groups: data.shiftMembers,
                  emptyMessage: '担当者はいません',
                ),
                const SizedBox(height: 12.0),
                _buildMemberSection(
                  title: '【前の時間の担当者の一覧】',
                  groups: [data.beforeMembers],
                  emptyMessage: '前の時間の担当者はいません',
                ),
                const SizedBox(height: 12.0),
                _buildMemberSection(
                  title: '【次の時間の担当者の一覧】',
                  groups: [data.afterMembers],
                  emptyMessage: '次の時間の担当者はいません',
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildMemberSection({
    required String title,
    required List<ShiftMembers> groups,
    required String emptyMessage,
  }) {
    final validGroups =
        groups.where((group) => group.members.isNotEmpty).toList();
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: const TextStyle(
            fontSize: AppFontSizes.xs,
            color: AppColors.textBlack,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 4.0),
        if (validGroups.isEmpty)
          Text(
            emptyMessage,
            style: const TextStyle(
              fontSize: AppFontSizes.xs,
              color: AppColors.textBlack,
              height: 1.5,
            ),
          )
        else
          ...validGroups.map(
            (group) => Padding(
              padding: const EdgeInsets.only(bottom: 4.0),
              child: Text(
                '${group.sTime}~${group.eTime}\n${group.members.map((member) => '(${member.bureau}${member.grade})${member.name}').join(', ')}',
                style: const TextStyle(
                  fontSize: AppFontSizes.xs,
                  color: AppColors.textBlack,
                  height: 1.4,
                ),
              ),
            ),
          ),
      ],
    );
  }
}

enum _ShiftCardMenuAction {
  members,
  review,
}

class _ShiftCardMoreButton extends StatelessWidget {
  static const double _menuLeftOffset = 0.0;
  static const double _menuTopOffset = 0.0;

  final Future<void> Function(Offset position) onPressed;

  const _ShiftCardMoreButton({required this.onPressed});

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: () => onPressed(_menuPosition(context)),
      tooltip: 'メニューを開く',
      padding: const EdgeInsets.all(8.0),
      constraints: const BoxConstraints(),
      style: IconButton.styleFrom(
        backgroundColor: Colors.transparent,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppBorderRadius.normal),
        ),
      ),
      icon: const Icon(
        Icons.more_vert,
        color: AppColors.textBlack,
        size: 24,
      ),
    );
  }

  Offset _menuPosition(BuildContext context) {
    final box = context.findRenderObject() as RenderBox;
    final topLeft = box.localToGlobal(Offset.zero);
    return Offset(
      topLeft.dx - _menuLeftOffset,
      topLeft.dy + _menuTopOffset,
    );
  }
}

class _ShiftCardMenuPanel extends PopupMenuEntry<_ShiftCardMenuAction> {
  const _ShiftCardMenuPanel();

  @override
  double get height => 138;

  @override
  bool represents(_ShiftCardMenuAction? value) => false;

  @override
  State<_ShiftCardMenuPanel> createState() => _ShiftCardMenuPanelState();
}

class _ShiftCardMenuPanelState extends State<_ShiftCardMenuPanel> {
  _ShiftCardMenuAction? _hoveredAction;
  _ShiftCardMenuAction? _pressedAction;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 1.0),
      child: DecoratedBox(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(8.0),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withValues(alpha: 0.40),
              offset: const Offset(0, 4),
              blurRadius: 4.0,
            ),
          ],
        ),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(8.0),
          child: ColoredBox(
            color: AppColors.base,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const SizedBox(height: 8.0),
                _ShiftCardMenuRow(
                  action: _ShiftCardMenuAction.members,
                  label: '担当者一覧',
                  isHovered: _hoveredAction == _ShiftCardMenuAction.members,
                  isPressed: _pressedAction == _ShiftCardMenuAction.members,
                  onHoverChanged: _handleHoverChanged,
                  onPressChanged: _handlePressChanged,
                ),
                _ShiftCardMenuRow(
                  action: _ShiftCardMenuAction.review,
                  label: 'レビューを書く',
                  isHovered: _hoveredAction == _ShiftCardMenuAction.review,
                  isPressed: _pressedAction == _ShiftCardMenuAction.review,
                  onHoverChanged: _handleHoverChanged,
                  onPressChanged: _handlePressChanged,
                ),
                const SizedBox(height: 8.0),
              ],
            ),
          ),
        ),
      ),
    );
  }

  void _handleHoverChanged(_ShiftCardMenuAction action, bool isHovered) {
    setState(() {
      _hoveredAction = isHovered ? action : null;
    });
  }

  void _handlePressChanged(_ShiftCardMenuAction action, bool isPressed) {
    setState(() {
      _pressedAction = isPressed ? action : null;
    });
  }
}

class _ShiftCardMenuRow extends StatelessWidget {
  final _ShiftCardMenuAction action;
  final String label;
  final bool isHovered;
  final bool isPressed;
  final void Function(_ShiftCardMenuAction action, bool isHovered)
      onHoverChanged;
  final void Function(_ShiftCardMenuAction action, bool isPressed)
      onPressChanged;

  const _ShiftCardMenuRow({
    required this.action,
    required this.label,
    required this.isHovered,
    required this.isPressed,
    required this.onHoverChanged,
    required this.onPressChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Semantics(
      button: true,
      label: label,
      child: InkWell(
        onHover: (hovering) => onHoverChanged(action, hovering),
        onTapDown: (_) => onPressChanged(action, true),
        onTapCancel: () => onPressChanged(action, false),
        onTap: () async {
          await Future<void>.delayed(const Duration(milliseconds: 120));
          if (context.mounted) {
            Navigator.pop(context, action);
          }
        },
        splashFactory: NoSplash.splashFactory,
        highlightColor: Colors.transparent,
        hoverColor: Colors.transparent,
        child: Container(
          width: double.infinity,
          height: 40,
          alignment: Alignment.center,
          color: _backgroundColor,
          child: Text(
            label,
            style: TextStyle(
              fontSize: AppFontSizes.sm,
              color: isPressed ? AppColors.textWhite : AppColors.textBlack,
            ),
          ),
        ),
      ),
    );
  }

  Color get _backgroundColor {
    if (isPressed) return AppColors.main;
    return AppColors.base;
  }
}

class _ManualToggle extends StatefulWidget {
  final String url;
  final VoidCallback? onOpened;

  const _ManualToggle({
    required this.url,
    this.onOpened,
  });

  @override
  State<_ManualToggle> createState() => _ManualToggleState();
}

class _ManualToggleState extends State<_ManualToggle> {
  bool _isExpanded = false;

  bool get _hasManual => widget.url.isNotEmpty;

  void _toggleManual() {
    if (!_hasManual) return;
    final willExpand = !_isExpanded;
    setState(() => _isExpanded = willExpand);
    if (willExpand && widget.onOpened != null) {
      widget.onOpened!();
    }
  }

  @override
  Widget build(BuildContext context) {
    final textColor = _hasManual ? AppColors.link : AppColors.grayDark;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        InkWell(
          borderRadius: BorderRadius.circular(4.0),
          onTap: _hasManual ? _toggleManual : null,
          child: Padding(
            padding: const EdgeInsets.symmetric(vertical: 4.0),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(
                  _isExpanded ? Icons.expand_less : Icons.expand_more,
                  size: 16,
                  color: textColor,
                ),
                const SizedBox(width: 4.0),
                Text(
                  _manualText,
                  style: TextStyle(
                    fontSize: AppFontSizes.xs,
                    color: textColor,
                    height: 1.5,
                  ),
                ),
              ],
            ),
          ),
        ),
        if (_isExpanded) ...[
          const SizedBox(height: 8.0),
          DecoratedBox(
            decoration: BoxDecoration(
              border: Border.all(color: AppColors.link),
              borderRadius: BorderRadius.circular(AppBorderRadius.normal),
            ),
            child: ClipRRect(
              borderRadius: BorderRadius.circular(AppBorderRadius.normal),
              child: Column(
                children: [
                  ManualViewer(url: widget.url),
                  Align(
                    alignment: Alignment.centerRight,
                    child: TextButton(
                      onPressed: () async {
                        await launchUrl(
                          Uri.parse(widget.url),
                          mode: LaunchMode.externalApplication,
                        );
                      },
                      child: const Text(
                        '別のタブで開く',
                        style: TextStyle(
                          fontSize: AppFontSizes.xs,
                          color: AppColors.link,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ],
    );
  }

  String get _manualText {
    if (!_hasManual) return 'マニュアルなし';
    return _isExpanded ? 'マニュアルを閉じる' : 'マニュアルを開く';
  }
}
