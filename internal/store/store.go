package store
import ("database/sql";"fmt";"os";"path/filepath";"strings";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Recording struct { ID string `json:"id"`; Title string `json:"title"`; Duration string `json:"duration,omitempty"`; Source string `json:"source,omitempty"`; Transcript string `json:"transcript,omitempty"`; Summary string `json:"summary,omitempty"`; Tags string `json:"tags,omitempty"`; Status string `json:"status"`; CreatedAt string `json:"created_at"`; WordCount int `json:"word_count"` }
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"scribe.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS recordings(id TEXT PRIMARY KEY,title TEXT NOT NULL,duration TEXT DEFAULT '',source TEXT DEFAULT '',transcript TEXT DEFAULT '',summary TEXT DEFAULT '',tags TEXT DEFAULT '',status TEXT DEFAULT 'pending',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(r *Recording)error{r.ID=genID();r.CreatedAt=now();if r.Status==""{r.Status="pending"};r.WordCount=len(strings.Fields(r.Transcript));_,err:=d.db.Exec(`INSERT INTO recordings VALUES(?,?,?,?,?,?,?,?,?)`,r.ID,r.Title,r.Duration,r.Source,r.Transcript,r.Summary,r.Tags,r.Status,r.CreatedAt);return err}
func(d *DB)Get(id string)*Recording{var r Recording;if d.db.QueryRow(`SELECT * FROM recordings WHERE id=?`,id).Scan(&r.ID,&r.Title,&r.Duration,&r.Source,&r.Transcript,&r.Summary,&r.Tags,&r.Status,&r.CreatedAt)!=nil{return nil};r.WordCount=len(strings.Fields(r.Transcript));return &r}
func(d *DB)List()[]Recording{rows,_:=d.db.Query(`SELECT * FROM recordings ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Recording;for rows.Next(){var r Recording;rows.Scan(&r.ID,&r.Title,&r.Duration,&r.Source,&r.Transcript,&r.Summary,&r.Tags,&r.Status,&r.CreatedAt);r.WordCount=len(strings.Fields(r.Transcript));o=append(o,r)};return o}
func(d *DB)Update(id string,r *Recording)error{r.WordCount=len(strings.Fields(r.Transcript));_,err:=d.db.Exec(`UPDATE recordings SET title=?,duration=?,source=?,transcript=?,summary=?,tags=?,status=? WHERE id=?`,r.Title,r.Duration,r.Source,r.Transcript,r.Summary,r.Tags,r.Status,id);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM recordings WHERE id=?`,id);return err}
func(d *DB)Search(q string)[]Recording{s:="%"+q+"%";rows,_:=d.db.Query(`SELECT * FROM recordings WHERE title LIKE ? OR transcript LIKE ? ORDER BY created_at DESC`,s,s);if rows==nil{return nil};defer rows.Close();var o []Recording;for rows.Next(){var r Recording;rows.Scan(&r.ID,&r.Title,&r.Duration,&r.Source,&r.Transcript,&r.Summary,&r.Tags,&r.Status,&r.CreatedAt);r.WordCount=len(strings.Fields(r.Transcript));o=append(o,r)};return o}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM recordings`).Scan(&n);return n}
