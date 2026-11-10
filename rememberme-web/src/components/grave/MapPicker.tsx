import { useEffect, useRef, useState, useCallback } from 'react'
import { Input, Button, Space, Spin, message, List, Typography } from 'antd'
import { AimOutlined, DeleteOutlined, LoadingOutlined, FullscreenOutlined, FullscreenExitOutlined } from '@ant-design/icons'
import AMapLoader from '@amap/amap-jsapi-loader'
import type { LatLng } from '../../types/grave'

const { Text } = Typography

const AMAP_KEY = import.meta.env.VITE_AMAP_KEY || ''
const AMAP_VERSION = '2.0'
const DEFAULT_CENTER: [number, number] = [116.397428, 39.90923]

let loadPromise: Promise<boolean> | null = null

function loadAMap(): Promise<boolean> {
  if (loadPromise) return loadPromise
  if ((window as any).AMap) {
    loadPromise = Promise.resolve(true)
    return loadPromise
  }
  if (!AMAP_KEY) {
    loadPromise = Promise.resolve(false)
    return loadPromise
  }

  ;(window as any)._AMapSecurityConfig = {
    securityJsCode: import.meta.env.VITE_AMAP_SECRET || '',
  }
  loadPromise = AMapLoader.load({
    key: AMAP_KEY,
    version: AMAP_VERSION,
    plugins: ['AMap.Geocoder', 'AMap.PlaceSearch', 'AMap.Geolocation'],
  })
    .then(() => true)
    .catch(() => false)

  return loadPromise
}

interface Props {
  value?: LatLng
  onChange: (v: LatLng | null) => void
}

export default function MapPicker({ value, onChange }: Props) {
  const containerRef = useRef<HTMLDivElement>(null)
  const mapRef = useRef<any>(null)
  const markerRef = useRef<any>(null)
  const geocoderRef = useRef<any>(null)
  const placeSearchRef = useRef<any>(null)
  const [loading, setLoading] = useState(true)
  const [mapError, setMapError] = useState(false)
  const [mapReady, setMapReady] = useState(false)
  const [searchText, setSearchText] = useState('')
  const [searchResults, setSearchResults] = useState<Array<{ name: string; address: string; lng: number; lat: number }>>([])
  const [searching, setSearching] = useState(false)
  const [manualLat, setManualLat] = useState(value ? String(value.latitude) : '')
  const [manualLng, setManualLng] = useState(value ? String(value.longitude) : '')
  const [showSearchList, setShowSearchList] = useState(false)
  const [isFullscreen, setIsFullscreen] = useState(false)

  const isInternalUpdate = useRef(false)
  const initializedRef = useRef(false)

  useEffect(() => {
    if (value) {
      setManualLat(String(value.latitude))
      setManualLng(String(value.longitude))
    } else {
      setManualLat('')
      setManualLng('')
    }

    if (!mapReady || !mapRef.current || isInternalUpdate.current) {
      isInternalUpdate.current = false
      return
    }
    if (value) {
      placeMarker(value.longitude, value.latitude)
    }
  }, [value, mapReady])

  const updatePosition = useCallback((lng: number, lat: number, address?: string) => {
    isInternalUpdate.current = true
    setManualLat(String(lat))
    setManualLng(String(lng))

    if (address !== undefined) {
      onChange({ latitude: lat, longitude: lng, address })
      return
    }

    if (geocoderRef.current) {
      geocoderRef.current.getAddress([lng, lat], (_status: string, result: any) => {
        const addr = (result?.regeocode?.formattedAddress) || ''
        onChange({ latitude: lat, longitude: lng, address: addr })
      })
    } else {
      onChange({ latitude: lat, longitude: lng, address: '' })
    }
  }, [onChange])

  const onSearch = useCallback((keyword: string) => {
    if (!keyword.trim() || !placeSearchRef.current) return
    setSearching(true)
    setShowSearchList(true)
    placeSearchRef.current.search(keyword, (_status: string, result: any) => {
      setSearching(false)
      console.log('[Amap PlaceSearch] status:', _status, 'result:', result)
      if (_status !== 'complete') {
        if (_status === 'no_data') {
          setSearchResults([])
        } else {
          message.warning('搜索失败，请检查高德 Key 和安全密钥配置')
        }
        return
      }
      if (result?.poiList?.pois) {
        setSearchResults(
          result.poiList.pois.map((poi: any) => ({
            name: poi.name,
            address: poi.address || '',
            lng: poi.location.lng,
            lat: poi.location.lat,
          }))
        )
      } else {
        setSearchResults([])
      }
    })
  }, [])

  const placeMarker = useCallback((lng: number, lat: number) => {
    const map = mapRef.current
    if (!map) return

    if (markerRef.current) {
      map.remove(markerRef.current)
    }

    const AMap = (window as any).AMap
    const marker = new AMap.Marker({
      position: [lng, lat],
      draggable: true,
    })
    map.add(marker)
    markerRef.current = marker
    map.setCenter([lng, lat])
    marker.on('dragend', () => {
      const pos = marker.getPosition()
      if (pos) updatePosition(pos.lng, pos.lat)
    })
    updatePosition(lng, lat)
  }, [updatePosition])

  const handleSearchSelect = useCallback((lng: number, lat: number) => {
    placeMarker(lng, lat)
    setShowSearchList(false)
  }, [placeMarker])

  const handleClear = useCallback(() => {
    if (markerRef.current && mapRef.current) {
      mapRef.current.remove(markerRef.current)
      markerRef.current = null
    }
    onChange(null)
    setManualLat('')
    setManualLng('')
  }, [onChange])

  const handleLocate = useCallback(() => {
    if (!navigator.geolocation) {
      message.warning('浏览器不支持定位')
      return
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => placeMarker(pos.coords.longitude, pos.coords.latitude),
      () => message.warning('获取定位失败，请检查定位权限'),
      { timeout: 10000 }
    )
  }, [placeMarker])

  const handleManualSubmit = useCallback(() => {
    const lat = parseFloat(manualLat)
    const lng = parseFloat(manualLng)
    if (isNaN(lat) || isNaN(lng)) {
      message.warning('请输入有效的经纬度')
      return
    }
    placeMarker(lng, lat)
  }, [manualLat, manualLng, placeMarker])

  useEffect(() => {
    setLoading(true)
    setMapError(false)
    loadPromise = null
    loadAMap().then((ok) => {
      if (!ok) {
        setLoading(false)
        if (AMAP_KEY) {
          setMapError(true)
        }
        return
      }
      const AMap = (window as any).AMap

      const map = new AMap.Map(containerRef.current, {
        zoom: 14,
        center: value ? [value.longitude, value.latitude] : DEFAULT_CENTER,
      })
      mapRef.current = map

      geocoderRef.current = new AMap.Geocoder({})
      placeSearchRef.current = new AMap.PlaceSearch({ pageSize: 5, pageIndex: 1 })

      if (value && value.latitude && value.longitude) {
        const marker = new AMap.Marker({
          position: [value.longitude, value.latitude],
          draggable: true,
        })
        map.add(marker)
        markerRef.current = marker
        marker.on('dragend', () => {
          const pos = marker.getPosition()
          if (pos) updatePosition(pos.lng, pos.lat)
        })
      }

      map.on('click', (e: any) => {
        placeMarker(e.lnglat.lng, e.lnglat.lat)
      })

      setMapReady(true)
      setLoading(false)
    })

    return () => {
      mapRef.current?.destroy()
      mapRef.current = null
      markerRef.current = null
      geocoderRef.current = null
      placeSearchRef.current = null
    }
  }, [])

  return (
    <div className={isFullscreen ? 'map-fullscreen-overlay' : ''}>
      {isFullscreen && (
        <div className="map-fs-bar">
          <Button size="small" icon={<FullscreenExitOutlined />} onClick={() => setIsFullscreen(false)}>
            退出全屏
          </Button>
        </div>
      )}
      <div className={isFullscreen ? 'map-fs-body' : ''}>
        <div style={{ position: 'relative', marginBottom: 8 }}>
          <Input.Search
            placeholder="搜索地点..."
            value={searchText}
            onChange={(e) => {
              setSearchText(e.target.value)
              if (!e.target.value.trim()) {
                setShowSearchList(false)
                setSearchResults([])
              }
            }}
            onSearch={onSearch}
            enterButton
            loading={searching}
          />
          {showSearchList && (
            <div
              style={{
                position: 'absolute',
                top: '100%',
                left: 0,
                right: 0,
                zIndex: 10,
                background: '#fff',
                border: '1px solid #d9d9d9',
                borderRadius: 4,
                maxHeight: 200,
                overflow: 'auto',
                boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
              }}
            >
              {searchResults.length === 0 && !searching ? (
                <div style={{ padding: 12, color: '#999', textAlign: 'center' }}>无结果</div>
              ) : (
                <List
                  size="small"
                  dataSource={searchResults}
                  renderItem={(item) => (
                    <List.Item
                      style={{ cursor: 'pointer', padding: '8px 12px' }}
                      onClick={() => handleSearchSelect(item.lng, item.lat)}
                    >
                      <div>
                        <div style={{ fontWeight: 500 }}>{item.name}</div>
                        <Text type="secondary" style={{ fontSize: 12 }}>{item.address}</Text>
                      </div>
                    </List.Item>
                  )}
                />
              )}
            </div>
          )}
        </div>

        <div style={{ position: 'relative', width: '100%', height: isFullscreen ? '100%' : 300, borderRadius: 8, border: '1px solid #d9d9d9', overflow: 'hidden' }}>
          <div ref={containerRef} style={{ width: '100%', height: '100%' }} />
          <Button
            className="map-fullscreen-btn"
            size="small"
            icon={<FullscreenOutlined />}
            onClick={() => setIsFullscreen(true)}
            style={{ display: isFullscreen ? 'none' : undefined }}
          />
          {loading && (
            <div style={{ position: 'absolute', inset: 0, display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#fafafa' }}>
              <Spin indicator={<LoadingOutlined style={{ fontSize: 24 }} />} tip="加载地图中..." />
            </div>
          )}
          {mapError && (
            <div style={{ position: 'absolute', inset: 0, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', background: '#fff7e6', color: '#999', fontSize: 13, gap: 8 }}>
              <span style={{ color: '#fa541c' }}>地图加载失败</span>
              <span style={{ fontSize: 12 }}>请检查网络连接，或刷新页面重试</span>
            </div>
          )}
          {!AMAP_KEY && !loading && (
            <div style={{ position: 'absolute', inset: 0, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', background: '#fafafa', color: '#999', fontSize: 13, gap: 8 }}>
              <span>高德地图 Key 未配置</span>
              <span style={{ fontSize: 12 }}>设置 <Text code>VITE_AMAP_KEY</Text> 后启用地图选点</span>
            </div>
          )}
        </div>

        {mapReady && (
          <div style={{ marginTop: 8 }}>
            <Space wrap>
              <Button size="small" icon={<AimOutlined />} onClick={handleLocate}>
                定位到当前位置
              </Button>
              <Button size="small" icon={<DeleteOutlined />} onClick={handleClear} disabled={!value?.latitude && !value?.longitude}>
                清除标记
              </Button>
              <Input
                size="small"
                placeholder="纬度"
                value={manualLat}
                onChange={(e) => setManualLat(e.target.value)}
                style={{ width: 110 }}
              />
              <Input
                size="small"
                placeholder="经度"
                value={manualLng}
                onChange={(e) => setManualLng(e.target.value)}
                style={{ width: 110 }}
              />
              <Button size="small" onClick={handleManualSubmit}>
                定位到坐标
              </Button>
            </Space>
          </div>
        )}

        {value?.address && (
          <div style={{ marginTop: 4, fontSize: 12, color: '#666' }}>
            {value.address}
          </div>
        )}
      </div>
    </div>
  )
}
