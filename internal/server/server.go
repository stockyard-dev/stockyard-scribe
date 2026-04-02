package server
import ("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-scribe/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux}
func New(db *store.DB)*Server{s:=&Server{db:db,mux:http.NewServeMux()}
s.mux.HandleFunc("GET /api/recordings",s.list);s.mux.HandleFunc("POST /api/recordings",s.create);s.mux.HandleFunc("GET /api/recordings/{id}",s.get);s.mux.HandleFunc("PUT /api/recordings/{id}",s.update);s.mux.HandleFunc("DELETE /api/recordings/{id}",s.del)
s.mux.HandleFunc("GET /api/search",s.search);s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"recordings":oe(s.db.List())})}
func(s *Server)create(w http.ResponseWriter,r *http.Request){var rec store.Recording;json.NewDecoder(r.Body).Decode(&rec);if rec.Title==""{we(w,400,"title required");return};s.db.Create(&rec);wj(w,201,s.db.Get(rec.ID))}
func(s *Server)get(w http.ResponseWriter,r *http.Request){rec:=s.db.Get(r.PathValue("id"));if rec==nil{we(w,404,"not found");return};wj(w,200,rec)}
func(s *Server)update(w http.ResponseWriter,r *http.Request){id:=r.PathValue("id");ex:=s.db.Get(id);if ex==nil{we(w,404,"not found");return};var rec store.Recording;json.NewDecoder(r.Body).Decode(&rec);if rec.Title==""{rec.Title=ex.Title};s.db.Update(id,&rec);wj(w,200,s.db.Get(id))}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)search(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"recordings":oe(s.db.Search(r.URL.Query().Get("q")))})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]int{"recordings":s.db.Count()})}
func(s *Server)health(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"status":"ok","service":"scribe","recordings":s.db.Count()})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
