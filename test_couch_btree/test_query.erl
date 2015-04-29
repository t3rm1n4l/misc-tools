-module(test_query).
-include("couch_db.hrl").
-export([start/0]).

calln(_Fn, 0) ->
    ok;
calln(Fn, N) ->
    Start = os:timestamp(),
    Fn(),
    io:format("time ~p~n", [timer:now_diff(os:timestamp(), Start)]),
    calln(Fn, N-1).

open_fast_fd(Filename, FilePid) ->
    {ok, Fd} = file2:open(Filename, [read, raw, binary]),
    erlang:put({FilePid, fast_fd_read}, Fd).

close_fast_fd(FilePid) ->
    Fd = erlang:erase({FilePid, fast_fd_read}),
    file:close(Fd).

loop(0) ->
    receive
    after 100000 ->
    ok
    end,
    ok;
loop(N) ->
    loop(N-1).



start() ->
    erlang:system_flag(scheduler_wall_time, true),
    lists:foreach(fun(_) -> spawn(fun() -> loop(1000000000) end) end, lists:seq(1,6000)),
    Filename = "/Users/sarath/development/couchbase/ns_server/data/n_0/data/@indexes/default/main_da1eaf6fac28abafd16daa38c3bbbfd7.view.1",
    couch_file_write_guard:sup_start_link(),
    {ok, Fd} = couch_file:open(Filename),
    Root = {41098790,<<0,0,15,66,63,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255>>,18982105},    %{ok, Btree}=couch_btree:open(NR, Fd),
    NewBtree0 = #btree{root=Root, fd=Fd, binary_mode=true},
    NewBtree = couch_btree:set_options(NewBtree0, [{less, less_fun}]),

    WrapperFn = fun() ->
    open_fast_fd(Filename, Fd),
    {_,_,Val}  = couch_btree:foldl(NewBtree, fun(KV, Red, X) ->
                                                 case X of
                                                 10 ->
                                                     {stop, X};
                                                 _ ->
                                                     {ok, X+1}
                                                 end
                                             end, 0),
    close_fast_fd(Fd),
    Val
    end,
    Ops = 550,
    Start = os:timestamp(),
    calln(WrapperFn, Ops),
    Stat=erlang:statistics(scheduler_wall_time),
    io:format("scheduler stats ~p~n", [Stat]),
    io:format("~p ops/secs~n", [Ops/(timer:now_diff(os:timestamp(), Start)/1000000)]).
