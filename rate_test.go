package rate

import (
	"testing"
	"time"
)

func TestLimiter_Allow(t *testing.T) {
	type fields struct {
		limit      int64
		windowSize time.Duration
		interval   time.Duration
		bucket     []int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "first",
			fields: fields{
				limit:      1,
				windowSize: 60 * time.Second,
				interval:   1 * time.Second,
				bucket:     make([]int64, 1),
			},
			want: true,
		},
		{
			name: "last",
			fields: fields{
				limit:      2,
				windowSize: 60 * time.Second,
				interval:   2 * time.Second,
				bucket:     []int64{time.Now().Add(-2 * time.Second).UnixNano()},
			},
			want: true,
		},
		{
			name: "interval not allowed",
			fields: fields{
				limit:      2,
				windowSize: 60 * time.Second,
				interval:   2 * time.Second,
				bucket:     []int64{time.Now().Add(-1 * time.Second).UnixNano()},
			},
			want: false,
		},
		{
			name: "windowSize not allowed",
			fields: fields{
				limit:      2,
				windowSize: 60 * time.Second,
				interval:   2 * time.Second,
				bucket:     []int64{time.Now().Add(-20 * time.Second).UnixNano(), time.Now().Add(-10 * time.Second).UnixNano()},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Limiter{
				limit:      tt.fields.limit,
				windowSize: tt.fields.windowSize,
				interval:   tt.fields.interval,
				bucket:     tt.fields.bucket,
			}
			if got := l.Allow(); got != tt.want {
				t.Errorf("Allow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimiter_Allow_1(t *testing.T) {
	rate := NewLimiter(3, 10*time.Second, 2*time.Second)
	after := time.After(15 * time.Second)
	for {
		select {
		case <-after:
			return
		default:
			if rate.Allow() {
				t.Logf("time:%s", time.Now().String())
			}
		}
	}
}
