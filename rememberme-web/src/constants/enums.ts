export const GRAVE_TYPE_MAP: Record<number, string> = {
  1: '单人墓',
  2: '双人墓',
  3: '家族墓',
}

export const GRAVE_TYPE_OPTIONS = [
  { value: 1, label: '单人墓' },
  { value: 2, label: '双人墓' },
  { value: 3, label: '家族墓' },
]

export const LAYOUT_MAP: Record<number, string> = {
  1: '单碑',
  2: '双碑并立',
  3: '多碑排列',
}

export const LAYOUT_OPTIONS = [
  { value: 1, label: '单碑' },
  { value: 2, label: '双碑并立' },
  { value: 3, label: '多碑排列' },
]

export const GENDER_MAP: Record<number, string> = {
  0: '未知',
  1: '男',
  2: '女',
}

export const GENDER_OPTIONS = [
  { value: 1, label: '男' },
  { value: 2, label: '女' },
]

export const GENDER_ICON: Record<number, string> = {
  1: '♂',
  2: '♀',
}

export const GRAVE_STATUS_OPTIONS = [
  { value: 1, label: '启用' },
  { value: 0, label: '停用' },
]

export const GRAVE_MODE_OPTIONS = [
  { value: 1, label: '纪念模式' },
  { value: 2, label: '族谱模式' },
]

export const RELATION_LABELS: Record<number, string> = {
  1: '儿子', 2: '女儿', 3: '配偶',
  4: '兄弟', 5: '姐妹', 6: '孙辈',
  7: '父亲', 8: '母亲', 9: '祖父', 10: '祖母',
  11: '长辈', 99: '追随者',
}

export const RELATION_OPTIONS = [
  { value: 1, label: '儿子' }, { value: 2, label: '女儿' }, { value: 3, label: '配偶' },
  { value: 4, label: '兄弟' }, { value: 5, label: '姐妹' }, { value: 6, label: '孙辈' },
  { value: 7, label: '父亲' }, { value: 8, label: '母亲' }, { value: 9, label: '祖父' },
  { value: 10, label: '祖母' }, { value: 11, label: '长辈' }, { value: 99, label: '追随者' },
]

export const DECEASED_RELATION_LABELS: Record<number, string> = {
  1: '夫妻',
  2: '亲子',
}

export const DECEASED_RELATION_OPTIONS = [
  { value: 1, label: '夫妻' },
  { value: 2, label: '亲子' },
]

export const RELATION_STATUS_MAP: Record<number, { label: string; color: string }> = {
  0: { label: '待确认', color: 'orange' },
  1: { label: '已确认', color: 'green' },
  2: { label: '已拒绝', color: 'red' },
}

export const ROLE_MAP: Record<number, { label: string; color: string }> = {
  10: { label: '所有者', color: 'gold' },
  20: { label: '管理员', color: 'blue' },
  30: { label: '成员', color: 'default' },
}

export const RECORD_TYPES: Record<number, { label: string; icon: string; color: string }> = {
  1: { label: '献花', icon: '🌸', color: 'green' },
  2: { label: '上香', icon: '🔥', color: 'orange' },
  3: { label: '点烛', icon: '🕯️', color: 'red' },
  4: { label: '鞠躬', icon: '🙏', color: 'blue' },
  5: { label: '留言', icon: '💬', color: 'purple' },
}

export const DIR_LABELS: Record<number, string> = {
  0: '对称',
  1: 'A 是长辈',
  2: 'B 是长辈',
}
