I wrote a SQL query to fetch spots with duplicate domains. Here's how it works:

I set up a join between two tables:

1.  (Table 1 - tb): This is the original table, but with an extra column called "website domain" that will be used for joining.

1.  (Table 2 - t2): It holds domains that appear more than once, along with their respective count.

To get the desired results, I joined "tb" with "t2" using a right outer join. It only shows rows that match the domains listed in "t2"'s website_domain column on "tb.". Now i can selectively choose the columns that i want to retrieved thanks to the join.
