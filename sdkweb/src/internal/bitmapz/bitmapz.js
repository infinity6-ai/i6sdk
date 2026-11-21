/**
 * @fileoverview Compact variable-length bitmask operations, serialization, and deserialization.
 * Direct JavaScript port of commonz/bitmapz.
 */

/**
 * Combines all given values into a single bitmask using bitwise OR.
 * @param {...(number|bigint)} values - Bitmask values to combine.
 * @returns {number} The combined bitmask.
 */
function join(...values) {
  let result = 0;
  for (const v of values) {
    result |= Number(v);
  }
  return result;
}

/**
 * Reports whether current contains all bits in expected.
 * If expected is empty, returns true.
 * @param {number|bigint} current - Bitmask to inspect.
 * @param {...(number|bigint)} expected - Required bits.
 * @returns {boolean} True if current contains all expected bits.
 */
function containsAll(current, ...expected) {
  const all = join(...expected);
  const cur = Number(current);
  return (cur & all) === all;
}

/**
 * Reports whether current contains at least one of the bits in expected.
 * If expected is empty or current is 0, returns false.
 * @param {number|bigint} current - Bitmask to inspect.
 * @param {...(number|bigint)} expected - Expected bits.
 * @returns {boolean} True if current contains any expected bit.
 */
function containsAny(current, ...expected) {
  const any = join(...expected);
  const cur = Number(current);
  return (cur & any) !== 0;
}

/**
 * Decomposes current into an array of individual set bits in ascending order.
 * If current is 0, returns an empty array.
 * @param {number|bigint} current - Bitmask to split.
 * @returns {number[]} Array of individual set bits.
 */
function split(current) {
  if (!current) return [];
  const result = [];
  let c = BigInt(current);
  while (c > 0n) {
    const bit = c & -c;
    result.push(Number(bit));
    c &= ~bit;
  }
  return result;
}

/**
 * Adds or removes the given values from current based on enable.
 * @param {number|bigint} current - Starting bitmask.
 * @param {boolean} enable - True to set bits, false to clear them.
 * @param {...(number|bigint)} values - Bitmask flags to add or remove.
 * @returns {number} The updated bitmask.
 */
function set(current, enable, ...values) {
  const mask = join(...values);
  const cur = Number(current);
  return enable ? (cur | mask) : (cur & ~mask);
}

/**
 * Encodes the bitmask v into a compact variable-length byte array (Uint8Array).
 * If v fits in 4 bits (0..15), it is serialized in 1 byte:
 * upper 4 bits are 0000, lower 4 bits store the value.
 * Otherwise, upper 4 bits specify the number of payload bytes (1..8),
 * lower 4 bits are 0000, followed by payload bytes in little-endian order.
 * @param {number|bigint} v - Bitmask value to serialize.
 * @returns {Uint8Array} The serialized bytes.
 */
function serialize(v) {
  const u = BigInt(v);
  if (u <= 0x0Fn) {
    return new Uint8Array([Number(u)]);
  }

  // Calculate bit length
  let temp = u;
  let bitLen = 0;
  while (temp > 0n) {
    bitLen++;
    temp >>= 1n;
  }

  const size = Math.floor((bitLen + 7) / 8);
  const buf = new Uint8Array(1 + size);
  buf[0] = size << 4;

  for (let i = 0; i < size; i++) {
    buf[1 + i] = Number((u >> BigInt(8 * i)) & 0xFFn);
  }
  return buf;
}

/**
 * Decodes a bitmask from data according to the variable-length encoding.
 * Throws an Error if data is empty, has an invalid length, invalid size,
 * or non-zero lower nibble when size > 0.
 * @param {Uint8Array} data - Binary encoded bitmask.
 * @returns {number|bigint} The decoded bitmask.
 */
function deserialize(data) {
  if (!data || data.length === 0) {
    throw new Error("bitmapz: data is empty");
  }

  const header = data[0];
  const size = header >> 4;

  if (size > 8) {
    throw new Error(`bitmapz: invalid size ${size}: bitmask size cannot exceed 8 bytes`);
  }

  if (size === 0) {
    if (data.length !== 1) {
      throw new Error(`bitmapz: invalid data length: expected 1 byte for 4-bit value, got ${data.length}`);
    }
    return header & 0x0F;
  }

  if ((header & 0x0F) !== 0) {
    throw new Error(`bitmapz: invalid header 0x${header.toString(16).padStart(2, '0')}: lower 4 bits must be 0 when size > 0`);
  }

  const expectedLen = 1 + size;
  if (data.length !== expectedLen) {
    throw new Error(`bitmapz: invalid data length: expected ${expectedLen} bytes for size ${size}, got ${data.length}`);
  }

  let val = 0n;
  for (let i = 0; i < size; i++) {
    val |= BigInt(data[1 + i]) << BigInt(8 * i);
  }

  if (val <= BigInt(Number.MAX_SAFE_INTEGER)) {
    return Number(val);
  }
  return val;
}

const bitmapz = Object.freeze({
  join,
  containsAll,
  containsAny,
  split,
  set,
  serialize,
  deserialize,
});

export { bitmapz, join, containsAll, containsAny, split, set, serialize, deserialize };
