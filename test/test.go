package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var mu sync.Mutex
var requestsDone int

func runClient(id, cycles int, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "127.0.0.1:64666")
	if err != nil {
		fmt.Println("client", id, "dial error:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024*8)

	for i := id * cycles; i < id*cycles+cycles; i++ {
		key := fmt.Sprintf("key%d", i)
		key_len := len(key)
		commands := []string{
			fmt.Sprintf("3,SET,%d,%s,4718,adsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtradsssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssu arehtucoshmrosexurhusodhmxguthxstuor,hxtuosr,hxtuiserhtuesrhothrethretcuoehructiweblutewruntuetbuiewrtbluxiebltuebrulbfdsuibcsobtprt'sbxdrtiobereobxw;erowxtberwlztnwerztkwemrtkwnctrowbctiowbxtubwrubywriotwcnrtm;opcentxioertuxbertuberuocthewriotierowthcrewipcthwuperwxtigtr", key_len, key),
			fmt.Sprintf("3,GET,%d,%s", key_len, key),
			fmt.Sprintf("3,SET,%d,%s,13,value changed", key_len, key),
			fmt.Sprintf("3,GET,%d,%s", key_len, key),
			fmt.Sprintf("3,DEL,%d,%s", key_len, key),
			fmt.Sprintf("3,GET,%d,%s", key_len, key),
		}

		for _, cmd := range commands {
			_, err := conn.Write([]byte(cmd + "\n"))
			if err != nil {
				fmt.Println("client", id, "write error:", err)
				return
			}
			_, err = conn.Read(buf)
			if err != nil {
				fmt.Println("client", id, "read error:", err)
				return
			}

			// Remove the comment below if you want to print
			//fmt.Printf("client %d: %s", id, string(buf[:n]))

			mu.Lock()
			requestsDone++
			mu.Unlock()
		}
	}
}

func main() {
	startedAt := time.Now()
	var wg sync.WaitGroup

	numClients := 25        // número de clientes concurrentes
	cyclesPerClient := 1000 // cada ciclo realiza 6 peticiones (SET, GET, SET, GET, DEL, GET)

	for c := 0; c < numClients; c++ {
		wg.Add(1)
		go runClient(c, cyclesPerClient, &wg)
	}

	wg.Wait()

	elapsed := time.Since(startedAt).Seconds()
	fmt.Println(requestsDone, "requests done!")
	fmt.Printf("-----------------\n%.2f seconds\n%.2f req/s\n",
		elapsed, float64(requestsDone)/elapsed)
}
