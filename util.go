package optparse

func appendString(arr []string, s string) []string {
    if len(arr) == cap(arr) {
        tmp := make([]string, len(arr), cap(arr) * 2);
        for i, e := range arr {
            tmp[i] = e;
        }
        arr = tmp;
    }
    arr = arr[0:len(arr)+1];
    arr[len(arr)-1] = s;
    return arr
}

func appendInt(arr []int, x int) []int {
    if len(arr) == cap(arr) {
        tmp := make([]int, len(arr), cap(arr) * 2);
        for i, e := range arr {
            tmp[i] = e;
        }
        arr = tmp;
    }
    arr = arr[0:len(arr)+1];
    arr[len(arr)-1] = x;
    return arr
}

func appendStringArray(arr [][]string, a []string) [][]string {
    if len(arr) == cap(arr) {
        tmp := make([][]string, len(arr), cap(arr) * 2);
        for i, e := range arr {
            tmp[i] = e;
        }
        arr = tmp;
    }
    arr = arr[0:len(arr)+1];
    arr[len(arr)-1] = a;
    return arr;
}
