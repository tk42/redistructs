package redistructs

// DeleteAll implements the types.Store interface.
// func (rs RedigoStructs) DeleteAll(ctx context.Context, mods ...rq.Modifier) error {
// 	conn, err := rs.pool.GetContext(ctx)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to acquire a connection")
// 	}
// 	defer conn.Close()

// 	dbIdx, changed := rs.getDBIndex(rs.model)
// 	if changed {
// 		conn.Do("SELECT", dbIdx)
// 	}

// 	keys, err := types.selectKeys(conn, mods)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}
// 	for _, k := range keys {
// 		_, err = conn.Do("HDEL", rm.name, k)
// 		if err != nil {
// 			return errors.Wrapf(err, "failed to remove by keys %v", keys)
// 		}

// 		_, err = conn.Do("ZREM", rs.name, k)
// 		if err != nil {
// 			return errors.Wrapf(err, "failed to remove by keys %v", keys)
// 		}
// 	}
// 	return nil
// }
