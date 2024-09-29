export class TestUtils {
  private static ZERO_PADDING_WIDTH = 6;

  static createPaddedID(prefix: string, value: number): string {
    return prefix + value.toString().padStart(this.ZERO_PADDING_WIDTH, '0');
  }

  static createRandomNumberGenerator(fourSeeds: number[]): () => number {
    // implements the SFC 32-bit PRNG from https://stackoverflow.com/a/47593316/650894

    function sfc32(a: number, b: number, c: number, d: number) {
      return function () {
        a |= 0; // invite JS runtime to do integer rather than float ops
        b |= 0;
        c |= 0;
        d |= 0;
        const t = (((a + b) | 0) + d) | 0;
        d = (d + 1) | 0;
        a = b ^ (b >>> 9);
        b = (c + (c << 3)) | 0;
        c = (c << 21) | (c >>> 11);
        c = (c + t) | 0;
        return (t >>> 0) / 4294967296; // in range [0, 1)
      };
    }
    return sfc32(fourSeeds[0], fourSeeds[1], fourSeeds[2], fourSeeds[3]);
  }
}
