import { Descriptions, Tag } from 'antd'
import { EnvironmentOutlined } from '@ant-design/icons'
import type { GraveResp } from '../../types/grave'
import { GRAVE_TYPE_MAP, LAYOUT_MAP } from '../../constants/enums'

interface Props {
  grave: GraveResp
}

export default function GraveDetail({ grave }: Props) {
  return (
    <Descriptions bordered column={{ xs: 1, sm: 2 }}>
      <Descriptions.Item label="墓园名称">{grave.name}</Descriptions.Item>
      <Descriptions.Item label="墓地类型">
        <Tag>{GRAVE_TYPE_MAP[grave.graveType] || '未知'}</Tag>
      </Descriptions.Item>
      <Descriptions.Item label="墓碑布局">
        <Tag>{LAYOUT_MAP[grave.layoutStyle] || '未知'}</Tag>
      </Descriptions.Item>
      <Descriptions.Item label="状态">
        {grave.status === 1 ? <Tag color="green">启用</Tag> : <Tag color="red">停用</Tag>}
      </Descriptions.Item>
      <Descriptions.Item label="创建者">{grave.creatorName}</Descriptions.Item>
      <Descriptions.Item label="创建时间">{grave.createdAt}</Descriptions.Item>
      <Descriptions.Item label="更新时间">{grave.updatedAt}</Descriptions.Item>

      <Descriptions.Item label="封面图片" span={2}>
        {grave.coverImage ? (
          <img src={grave.coverImage} alt="封面" style={{ maxWidth: 300, maxHeight: 200 }} />
        ) : (
          '暂无封面'
        )}
      </Descriptions.Item>

            <Descriptions.Item label="墓园位置" span={2}>
        {grave.latitude != null && grave.longitude != null ? (
          <span>
            <EnvironmentOutlined style={{ marginRight: 4, color: '#1890ff' }} />
            {grave.address || `${grave.latitude.toFixed(6)}, ${grave.longitude.toFixed(6)}`}
            <span style={{ marginLeft: 12, color: '#999', fontSize: 12 }}>
              ({grave.latitude.toFixed(6)}, {grave.longitude.toFixed(6)})
            </span>
          </span>
        ) : (
          '未设置'
        )}
      </Descriptions.Item>

      <Descriptions.Item label="描述" span={2}>
        {grave.description || '暂无描述'}
      </Descriptions.Item>
    </Descriptions>
  )
}
