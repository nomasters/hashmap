package hashmap

// // Engine is the enum type for StorageEngine
// type Engine int

// // Enum types for Storage Engine
// const (
// 	MemoryEngine Engine = iota
// 	RedisEngine
// )

// // StorageOptions are used to bootstrap a new StorageEngine. Currently this is intended for Redis, but could be
// // expanded to other storage types in the future.
// type StorageOptions struct {
// 	Engine          Engine
// 	Endpoint        string
// 	Auth            string
// 	MaxIdle         int
// 	MaxActive       int
// 	IdleTimeout     time.Duration
// 	Wait            bool
// 	MaxConnLifetime time.Duration
// 	TLS             bool
// }

// // Storage Defaults
// const (
// 	MetadataPrefix      = "meta-"
// 	DefaultRedisAddress = ":6379"
// )

// // String is used to pretty print storage engine constants
// func (s Engine) String() string {
// 	names := []string{
// 		"memory",
// 		"redis",
// 	}
// 	return names[s]
// }

// // GetStorageEngineCode takes a storage engine name in the form of a string
// // and returns and Engine or an error
// func GetStorageEngineCode(n string) (Engine, error) {
// 	switch n {
// 	case "memory":
// 		return MemoryStorage, nil
// 	case "redis":
// 		return RedisStorage, nil
// 	default:
// 		return 0, errors.New("invalid storage engine name")
// 	}
// }

// // NewStorage is a helper function used for configuring supported storage engines
// func NewStorage(opts StorageOptions) (Storage, error) {
// 	switch opts.Engine {
// 	case MemoryStorage:
// 		return NewMemoryStore(), nil
// 	case RedisStorage:
// 		return NewRedisStore(opts), nil
// 	default:
// 		return nil, errors.New("invalid storage engine")
// 	}
// }

// // Storage is the primary interface for interacting with Payload and PayloadMetaData
// type Storage interface {
// 	Get(key string) (PayloadWithMetadata, error)
// 	Set(key string, value PayloadWithMetadata) error
// 	Delete(key string) error
// }

// // MemoryStore  is the primary in-memory data storage and retrieval struct
// type MemoryStore struct {
// 	sync.RWMutex
// 	internal map[string]PayloadWithMetadata
// }

// // NewMemoryStore returns a pointer to a new intance ofMemoryStore
// func NewMemoryStore() *MemoryStore {
// 	return &MemoryStore{
// 		internal: make(map[string]PayloadWithMetadata),
// 	}
// }

// // Get method for MemoryStore with read locks
// func (m *MemoryStore) Get(key string) (PayloadWithMetadata, error) {
// 	var err error
// 	m.RLock()
// 	v, ok := m.internal[key]
// 	m.RUnlock()
// 	if !ok {
// 		err = errors.New("key not found")
// 	}
// 	return v, err
// }

// // Set method for MemoryStore with read/write locks
// func (m *MemoryStore) Set(key string, value PayloadWithMetadata) error {
// 	m.Lock()
// 	m.internal[key] = value
// 	m.Unlock()
// 	return nil
// }

// // Delete method for MemoryStore with read/write locks
// func (m *MemoryStore) Delete(key string) error {
// 	m.Lock()
// 	delete(m.internal, key)
// 	m.Unlock()
// 	return nil
// }

// // RedisStore is a struct with methods that conforms to the Storage Interface
// type RedisStore struct {
// 	pool *redis.Pool
// }

// // NewRedisStore returns a RedisStore with StorageOptions mapped to Redis Pool settings.
// // Optionally, if Auth is set, Auth is configured on Dial
// func NewRedisStore(opts StorageOptions) *RedisStore {
// 	addr := opts.Endpoint
// 	if addr == "" {
// 		addr = DefaultRedisAddress
// 	}

// 	return &RedisStore{
// 		pool: &redis.Pool{
// 			MaxIdle:         opts.MaxIdle,
// 			MaxActive:       opts.MaxActive,
// 			IdleTimeout:     opts.IdleTimeout,
// 			Wait:            opts.Wait,
// 			MaxConnLifetime: opts.MaxConnLifetime,
// 			Dial: func() (redis.Conn, error) {
// 				// TODO: DialTLSSkipVerify needs to be made to an optional param,
// 				// this is a hack put in place because AWS ElasticCache failed to x509 verify
// 				c, err := redis.Dial("tcp", addr, redis.DialUseTLS(opts.TLS), redis.DialTLSSkipVerify(true))
// 				if err != nil {
// 					return nil, err
// 				}
// 				if opts.Auth != "" {
// 					if _, err := c.Do("AUTH", opts.Auth); err != nil {
// 						c.Close()
// 						return nil, err
// 					}
// 				}
// 				return c, nil
// 			},
// 		},
// 	}
// }

// // Get method for RedisStore
// func (r *RedisStore) Get(key string) (PayloadWithMetadata, error) {
// 	var pwm PayloadWithMetadata
// 	c := r.pool.Get()
// 	defer c.Close()

// 	response, err := redis.StringMap(c.Do("HGETALL", key))
// 	if err != nil {
// 		return pwm, err
// 	}

// 	log.Println(response)

// 	mp := []byte(response["payload"])
// 	p := Payload{}
// 	if err := json.Unmarshal(mp, &p); err != nil {
// 		return pwm, err
// 	}

// 	pwm.Payload = p

// 	for k, v := range response {
// 		if strings.HasPrefix(k, MetadataPrefix) {
// 			pwm.Metadata[k[len(MetadataPrefix):]] = v
// 		}
// 	}

// 	return pwm, nil
// }

// // Set method of RedisStore
// func (r *RedisStore) Set(key string, value PayloadWithMetadata) error {

// 	data, err := value.Payload.GetData()
// 	if err != nil {
// 		return err
// 	}

// 	ttl := data.TTL
// 	mp, err := json.Marshal(value.Payload)
// 	if err != nil {
// 		return err
// 	}
// 	c := r.pool.Get()
// 	defer c.Close()
// 	// atomic blog for writing all hash values and then setting TTL
// 	c.Send("MULTI")
// 	c.Send("HSET", key, "payload", string(mp))

// 	// add prefix to hash key to make it easy to parse
// 	for k, v := range value.Metadata {
// 		mk := MetadataPrefix + k
// 		c.Send("HSET", key, mk, v)
// 	}

// 	c.Send("EXPIRE", key, ttl)
// 	if _, err := c.Do("EXEC"); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Delete method for RedisStore
// func (r *RedisStore) Delete(key string) error {
// 	c := r.pool.Get()
// 	defer c.Close()

// 	if _, err := c.Do("DEL", key); err != nil {
// 		return err
// 	}

// 	return nil
// }
